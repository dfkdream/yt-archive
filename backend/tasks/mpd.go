package tasks

import (
	"errors"
	"log/slog"
	"strconv"

	"github.com/zencoder/go-dash/mpd"
)

var (
	ErrNoPeriod = errors.New("no period found")
)

func mergeMPDs(dstPath string, srcPaths ...string) error {
	slog.Info("merging mpds.", "dst", dstPath, "srcs", srcPaths)
	d, err := mpd.ReadFromFile(dstPath)
	if err != nil {
		return err
	}

	for _, p := range srcPaths {
		s, err := mpd.ReadFromFile(p)
		if err != nil {
			return err
		}
		mergeMPD(d, s)
	}

	return d.WriteToFile(dstPath)
}

func mergeMPD(dst *mpd.MPD, src *mpd.MPD) error {
	if len(src.Periods) == 0 || len(dst.Periods) == 0 {
		return ErrNoPeriod
	}

	var videos *mpd.AdaptationSet
	var audios *mpd.AdaptationSet

	for _, a := range dst.Periods[0].AdaptationSets {
		if *a.MimeType == "video/webm" {
			videos = a
		} else {
			audios = a
		}
	}

	var srcPeriod = src.Periods[0]
	var dstPeriod = dst.Periods[0]

	for _, a := range srcPeriod.AdaptationSets {
		if *a.MimeType == "video/webm" {
			if videos == nil {
				dstPeriod.AdaptationSets = append(srcPeriod.AdaptationSets, dstPeriod.AdaptationSets...)
				break
			} else {
				needUpdate, err := isResUpdateRequired(videos, a)
				if err != nil {
					return err
				}
				if needUpdate {
					slog.Info("updating adaptation set resolution")
					videos.Width = a.Width
					videos.Height = a.Height
				}
				videos.Representations = append(videos.Representations, a.Representations...)
			}
		} else {
			if audios == nil {
				dstPeriod.AdaptationSets = append(dstPeriod.AdaptationSets, srcPeriod.AdaptationSets...)
				break
			} else {
				audios.Representations = append(audios.Representations, a.Representations...)
			}
		}
	}

	repID := 0
	for i, a := range dstPeriod.AdaptationSets {
		*a.ID = strconv.Itoa(i)
		for _, r := range a.Representations {
			*r.ID = strconv.Itoa(repID)
			repID++
		}
	}

	return nil
}

func isResUpdateRequired(src, target *mpd.AdaptationSet) (bool, error) {
	if target.Width == nil || target.Height == nil {
		return false, nil
	}

	if src.Width == nil || src.Height == nil {
		return true, nil
	}

	sw, err := strconv.Atoi(*src.Width)
	if err != nil {
		return false, err
	}

	sh, err := strconv.Atoi(*src.Height)
	if err != nil {
		return false, err
	}

	tw, err := strconv.Atoi(*target.Width)
	if err != nil {
		return false, err
	}

	th, err := strconv.Atoi(*target.Height)
	if err != nil {
		return false, err
	}

	return tw > sw && th > sh, nil
}
