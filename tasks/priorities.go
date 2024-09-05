package tasks

const (
	PRIORITY_GAP = 20
)

const (
	PRIORITY_LOWEST = iota * PRIORITY_GAP
	PRIORITY_DOWNLOAD_VIDEO
	PRIORITY_DOWNLOAD_AUDIO
	PRIORITY_ARCHIVE_PLAYLIST
	PRIORITY_ARCHIVE_VIDEO
	PRIORITY_ARCHVIE_CHANNEL_INFO
	PRIORITY_HIGHEST
)

func calculateVideoPriority(f format) int {
	priority := PRIORITY_DOWNLOAD_VIDEO

	if f.VideoExt == "webm" {
		priority += PRIORITY_GAP / 2
	}

	pixels := f.Width * f.Height

	switch {
	case pixels <= 256*144:
		priority += 6
	case pixels <= 426*240:
		priority += 5
	case pixels <= 640*360:
		priority += 4
	case pixels <= 854*480:
		priority += 3
	case pixels <= 1280*720:
		priority += 2
	case pixels <= 1920*1080:
		priority += 1
	}

	return priority
}
