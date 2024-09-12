package tasks

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"yt-archive/taskq"
)

const TaskTypeArchiveVideo = "ARCHIVE_VIDEO"

type ArchiveVideoHandler struct {
	DB *sql.DB
}

func NewArchiveVideoHandler(db *sql.DB) (ArchiveVideoHandler, error) {
	_, err := db.Exec("create table if not exists videos (id text primary key, title text, description text, timestamp timestamp, duration text, owner text, thumbnail text)")
	if err != nil {
		return ArchiveVideoHandler{}, err
	}

	return ArchiveVideoHandler{DB: db}, nil
}

type format struct {
	FormatID   string  `json:"format_id"`
	VideoExt   string  `json:"video_ext"`
	VideoCodec string  `json:"vcodec"`
	AudioCodec string  `json:"acodec"`
	FPS        float32 `json:"fps"`
	Width      int     `json:"width"`
	Height     int     `json:"height"`
	VBR        float32 `json:"vbr"`
	Protocol   string  `json:"protocol"`
	Name       string  `json:"format"`
}

type videoMetadata struct {
	ID          string   `json:"id"`
	Title       string   `json:"fulltitle"`
	Description string   `json:"description"`
	Timestamp   int      `json:"timestamp"`
	Duration    string   `json:"duration_string"`
	Owner       string   `json:"uploader_id"`
	Formats     []format `json:"formats"`
}

func (a ArchiveVideoHandler) Handler(task *taskq.Task) error {
	var videoID string
	err := json.Unmarshal(task.Payload, &videoID)
	if err != nil {
		return err
	}

	r := a.DB.QueryRow("select count(id) from videos where id=?", videoID)
	var n int
	err = r.Scan(&n)
	if err != nil {
		return err
	}

	if n > 0 {
		slog.Info("skipping duplicated video", "id", videoID)
		return nil
	}

	tempDir, err := os.MkdirTemp("", videoID+"_*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	err = Exec("yt-dlp", "--write-info-json", "--skip-download", "--write-thumbnail", "-o", "%(id)s.%(ext)s", "--paths", tempDir, videoID)
	if err != nil {
		return err
	}

	destPath := filepath.Join("videos", videoID)
	err = os.MkdirAll(destPath, os.FileMode(0o700))
	if err != nil {
		return err
	}

	thumbnail, err := copyThumbnail(tempDir, destPath)
	if err != nil {
		return err
	}

	slog.Info("downloaded thumbnail", "filename", thumbnail)

	metadata, err := parseVideoMetadata(filepath.Join(tempDir, videoID+".info.json"))
	if err != nil {
		return err
	}

	_, err = a.DB.Exec("insert into videos (id, title, description, timestamp, duration, owner, thumbnail) values (?, ?, ?, ?, ?, ?, ?)",
		metadata.ID, metadata.Title, metadata.Description, time.Unix(int64(metadata.Timestamp), 0), metadata.Duration, metadata.Owner, thumbnail)
	if err != nil {
		return err
	}

	t, err := taskq.NewJsonTask(PriorityArchiveChannelInfo, TaskTypeArchiveChannelInfo, videoID+"_"+metadata.Owner, metadata.Owner)
	if err != nil {
		return err
	}

	taskq.Enqueue(t)

	formats := selectVideoFormats(metadata.Formats)

	downloadMediaPayload := DownloadMediaPayload{
		VideoID:      videoID,
		Format:       "bestaudio",
		OutputPath:   filepath.Join(destPath, videoID+AUDIO_FILE_SUFFIX),
		SkipEncoding: false,
	}

	t, err = taskq.NewJsonTask(PriorityDownloadAudio, TaskTypeDownloadMedia, videoID+"_bestaudio", downloadMediaPayload)
	if err != nil {
		return err
	}

	taskq.Enqueue(t)

	for _, v := range formats {
		downloadMediaPayload.Format = v.FormatID
		downloadMediaPayload.OutputPath = filepath.Join(destPath, videoID+"_"+strconv.Itoa(v.Height)+MEDIA_FILE_SUFFIX)
		downloadMediaPayload.SkipEncoding = canSkipEncoding(v)

		description := fmt.Sprintf("%s, %d, %s", videoID, v.Height, v.VideoCodec)
		if !downloadMediaPayload.SkipEncoding {
			description += " (Encoding required)"
		}

		t, err = taskq.NewJsonTask(calculateVideoPriority(v), TaskTypeDownloadMedia, description, downloadMediaPayload)
		if err != nil {
			return err
		}

		taskq.Enqueue(t)
	}

	return nil
}

func selectVideoFormats(formats []format) []format {
	/*
		Format preference:
		webm (vp09, av01) - Other

		Protocol preference:
		https - Other

		Resolution preference:
		144 - 240 - 360 - 480 - 720 - 1080 (Sort by VBR)
	*/

	result := make([]format, 0)

	resolutionFormat := make(map[string]format)

	for _, f := range formats {
		if f.VideoCodec == "none" {
			continue
		}

		if f.AudioCodec != "none" {
			// skip muxed files
			continue
		}

		resolutionString := fmt.Sprintf("%dx%d", f.Width, f.Height)
		prev, ok := resolutionFormat[resolutionString]
		if !ok {
			resolutionFormat[resolutionString] = f
			continue
		}

		if calculateFormatPreference(f) > calculateFormatPreference(prev) {
			resolutionFormat[resolutionString] = f
		}
	}

	for _, v := range resolutionFormat {
		result = append(result, v)
	}

	return result
}

func calculateFormatPreference(f format) int {
	preference := 0

	if f.Protocol == "https" {
		// DASH
		preference++
	}

	if canSkipEncoding(f) {
		preference++
	}

	if f.FPS > 30 {
		preference++
	}

	return preference
}

func canSkipEncoding(f format) bool {
	if strings.HasPrefix(f.VideoCodec, "vp") {
		// VP8 or VP9
		return true
	}

	if strings.HasPrefix(f.VideoCodec, "av01") {
		// AV1
		return true
	}

	return false
}

func parseVideoMetadata(path string) (*videoMetadata, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var metadata videoMetadata
	err = json.NewDecoder(f).Decode(&metadata)
	return &metadata, err
}

const thumbnailExtensions = ".webp|.jpg|.png"

func copyThumbnail(path string, dest string) (string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}

	filename := ""
	for _, f := range files {
		ext := filepath.Ext(f.Name())
		if strings.Contains(thumbnailExtensions, ext) {
			filename = f.Name()
		}
	}

	if filename == "" {
		return "", errors.New("could not find thumbnail in directory " + path)
	}

	src, err := os.Open(filepath.Join(path, filename))
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(filepath.Join(dest, filename))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)

	return filename, err
}
