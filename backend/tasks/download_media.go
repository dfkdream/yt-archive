package tasks

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"yt-archive/taskq"
)

const TaskTypeDownloadMedia = "DOWNLOAD_MEDIA"

const (
	MEDIA_FILE_SUFFIX = ".webm"
	AUDIO_FILE_SUFFIX = "_audio" + MEDIA_FILE_SUFFIX
)

type DownloadMediaPayload struct {
	VideoID      string
	Format       string
	OutputPath   string
	SkipEncoding bool
}

func DownloadMediaHandler(task *taskq.Task) error {
	var payload DownloadMediaPayload
	err := json.Unmarshal(task.Payload, &payload)
	if err != nil {
		return err
	}

	tempDir, err := os.MkdirTemp("", payload.VideoID+"_*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	videoURL := "https://www.youtube.com/watch?v=" + payload.VideoID
	err = Exec("yt-dlp", "-f", payload.Format, "-o", "%(id)s.%(ext)s", "--paths", tempDir, videoURL)
	if err != nil {
		return err
	}

	files, err := os.ReadDir(tempDir)
	if err != nil {
		return err
	}

	if len(files) != 1 {
		return errors.New("expected 1 file in directory but got " + strconv.Itoa(len(files)))
	}

	filename := files[0].Name()
	slog.Info("download completed", "filename", filename)

	if (payload.Format == "bestaudio" && filepath.Ext(filename) == ".webm") || payload.SkipEncoding {
		err = Exec("ffmpeg", "-hide_banner", "-y", "-i", filepath.Join(tempDir, filename), "-f", "webm", "-c", "copy", "-dash", "1", payload.OutputPath)
	} else {
		err = Exec("ffmpeg", "-hide_banner", "-y", "-i", filepath.Join(tempDir, filename), "-keyint_min", "150", "-g", "150", "-tile-columns", "4", "-frame-parallel", "1", "-f", "webm", "-dash", "1", payload.OutputPath)
	}

	if err != nil {
		return err
	}

	tempMpd := filepath.Join(tempDir, "manifest.mpd")
	err = buildManifest(payload.OutputPath, tempMpd)
	if err != nil {
		return err
	}

	masterMpd := filepath.Join(payload.OutputPath, "..", payload.VideoID+".mpd")
	return mergeMPDs(masterMpd, tempMpd)
}

func buildManifest(path string, output string) error {
	command := []string{
		"-hide_banner", "-y",
		"-f", "webm_dash_manifest", "-i", path,
		"-c", "copy",
		"-map", "0",
		"-f", "webm_dash_manifest",
		"-adaptation_sets", "id=0,streams=0", output,
	}

	return Exec("ffmpeg", command...)
}
