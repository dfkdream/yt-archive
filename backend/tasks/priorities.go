package tasks

const (
	priorityGap = 20
)

const (
	PriorityLowest = iota * priorityGap
	PriorityDownloadVideo
	PriorityDownloadAudio
	PriorityArchivePlaylist
	PriorityArchiveVideo
	PriorityArchiveChannelInfo
	PriorityHighest
)

func CalculateVideoPriority(f format) int {
	priority := PriorityDownloadVideo

	if CanSkipEncoding(f) {
		priority += priorityGap / 2
	}

	pixels := f.Width * f.Height

	switch {
	case pixels <= 256*144:
		priority += 7
	case pixels <= 426*240:
		priority += 6
	case pixels <= 640*360:
		priority += 5
	case pixels <= 854*480:
		priority += 4
	case pixels <= 1280*720:
		priority += 3
	case pixels <= 1920*1080:
		priority += 2
	case pixels <= 2560*1440:
		priority += 1
	}

	return priority
}
