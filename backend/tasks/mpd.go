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
				needUpdate, err := isHeigherRes(a, videos)
				if err != nil {
					return err
				}
				if needUpdate {
					*videos.Width = *a.Width
					*videos.Height = *a.Height
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

func isHeigherRes(left, right *mpd.AdaptationSet) (bool, error) {
	lw, err := strconv.Atoi(*left.Width)
	if err != nil {
		return false, err
	}

	lh, err := strconv.Atoi(*left.Height)
	if err != nil {
		return false, err
	}

	rw, err := strconv.Atoi(*right.Width)
	if err != nil {
		return false, err
	}

	rh, err := strconv.Atoi(*right.Height)
	if err != nil {
		return false, err
	}

	return lw > rw && lh > rh, nil
}
