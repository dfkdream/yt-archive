package tasks

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

	outputDir := filepath.Join(payload.OutputPath, "..")
	manifestPath := filepath.Join(outputDir, payload.VideoID+".mpd")
	return BuildManifest(outputDir, manifestPath)
}

func BuildManifest(path string, output string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	videoFiles := make([]string, 0)
	audioFile := ""
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		if filepath.Ext(f.Name()) != MEDIA_FILE_SUFFIX {
			continue
		}

		filepath := filepath.Join(path, f.Name())

		if strings.HasSuffix(f.Name(), AUDIO_FILE_SUFFIX) {
			audioFile = filepath
		} else {
			videoFiles = append(videoFiles, filepath)
		}
	}

	command := []string{"-hide_banner", "-y"}

	for _, filepath := range videoFiles {
		command = append(command, "-f", "webm_dash_manifest", "-i", filepath)
	}

	if audioFile != "" {
		command = append(command, "-f", "webm_dash_manifest", "-i", audioFile)
	}

	command = append(command, "-c", "copy")

	filecount := len(videoFiles)
	if audioFile != "" {
		filecount += 1
	}

	for i := 0; i < filecount; i++ {
		command = append(command, "-map", strconv.Itoa(i))
	}

	command = append(command, "-f", "webm_dash_manifest")

	adaptationSets := "id=0,streams="

	for i := range videoFiles {
		adaptationSets += strconv.Itoa(i)

		if i < len(videoFiles)-1 {
			adaptationSets += ","
		} else if audioFile != "" {
			adaptationSets += " id=1,streams="
		}
	}

	if audioFile != "" {
		adaptationSets += strconv.Itoa(len(videoFiles))
	}

	command = append(command, "-adaptation_sets", adaptationSets, output)

	return Exec("ffmpeg", command...)
}
