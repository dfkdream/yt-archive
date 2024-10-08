package mpd

func Merge(m1, m2 MPD) MPD {
	result := m1

	if m2.MediaPresentationDuration > m1.MediaPresentationDuration {
		result.MediaPresentationDuration = m2.MediaPresentationDuration
	}

	var videoSet AdaptationSet
	var audioSet AdaptationSet

	for _, a := range append(m1.Period.AdaptationSet, m2.Period.AdaptationSet...) {
		if a.Width != 0 && a.Height != 0 && len(a.Representation) == 1 {
			a.Representation[0].Width = a.Width
			a.Representation[0].Height = a.Height
			a.Width = 0
			a.Height = 0
		}

		switch a.MimeType {
		case "video/webm":
			if videoSet.Representation == nil {
				videoSet = a
			} else {
				if videoSet.Lang == "" && a.Lang != "" {
					videoSet.Lang = a.Lang
				}
				videoSet.SubsegmentAlignment = false
				videoSet.Representation = append(videoSet.Representation, a.Representation...)
			}
		case "audio/webm":
			if audioSet.Representation == nil {
				audioSet = a
			} else {
				if audioSet.Lang == "" && a.Lang != "" {
					audioSet.Lang = a.Lang
				}
				audioSet.SubsegmentAlignment = false
				audioSet.Representation = append(audioSet.Representation, a.Representation...)
			}
		}
	}

	adaptationSet := []AdaptationSet{}

	if videoSet.Representation != nil {
		adaptationSet = append(adaptationSet, videoSet)
	}

	if audioSet.Representation != nil {
		adaptationSet = append(adaptationSet, audioSet)
	}

	if len(adaptationSet) > 0 {
		result.Period.AdaptationSet = adaptationSet
	}

	reassignIDs(&result)

	return result
}

func reassignIDs(mpd *MPD) {
	repID := 0
	for i := range mpd.Period.AdaptationSet {
		a := &mpd.Period.AdaptationSet[i]
		a.Id = i
		for j := range a.Representation {
			a.Representation[j].Id = repID
			repID++
		}
	}
}
