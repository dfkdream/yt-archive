package tasks

import (
	"errors"
	"io/fs"
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

	if errors.Is(err, fs.ErrNotExist) {
		d = nil
	} else if err != nil {
		return err
	}

	for _, p := range srcPaths {
		s, err := mpd.ReadFromFile(p)
		if err != nil {
			return err
		}

		if d == nil {
			d = s
			s = nil
		}

		mergeMPD(d, s)
	}

	return d.WriteToFile(dstPath)
}

func mergeMPD(dst *mpd.MPD, src *mpd.MPD) error {
	if src == nil {
		for _, a := range dst.Periods[0].AdaptationSets {
			if *a.MimeType == "video/webm" && a.Width != nil && a.Height != nil && len(a.Representations) == 1 {
				return updateSingleVideoAdaptationSet(a)
			}
		}

		return nil
	}

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
			if a.Width != nil && a.Height != nil && len(a.Representations) == 1 {
				// single video adaptationSet
				slog.Info("single video adaptation set found. updating representation.")
				updateSingleVideoAdaptationSet(a)
			}
		}
	}

	for _, a := range srcPeriod.AdaptationSets {
		if *a.MimeType == "video/webm" {
			if videos == nil {
				dstPeriod.AdaptationSets = append(srcPeriod.AdaptationSets, dstPeriod.AdaptationSets...)
				break
			} else {
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

func updateSingleVideoAdaptationSet(a *mpd.AdaptationSet) error {
	width, err := strconv.ParseInt(*a.Width, 10, 64)
	if err != nil {
		return err
	}

	height, err := strconv.ParseInt(*a.Height, 10, 64)
	if err != nil {
		return err
	}

	a.Width = nil
	a.Height = nil
	a.Representations[0].Width = &width
	a.Representations[0].Height = &height

	return nil
}
