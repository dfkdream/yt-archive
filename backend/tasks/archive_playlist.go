package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
	"yt-archive/db"
	"yt-archive/taskq"
)

const TaskTypeArchivePlaylist = "ARCHIVE_PLAYLIST"

type playlistMetadata struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Owner       string `json:"uploader_id"`
	Timestamp   string `json:"modified_date"`
}

func ArchivePlaylistHandler(task *taskq.Task) error {
	var playlistID string
	err := json.Unmarshal(task.Payload, &playlistID)
	if err != nil {
		return err
	}

	tempDir, err := os.MkdirTemp("", playlistID+"_*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	err = Exec("yt-dlp", "--write-info-json", "--skip-download", "-o", "%(id)s.%(ext)s", "--paths", tempDir, playlistID)
	if err != nil {
		return err
	}

	files, err := os.ReadDir(tempDir)
	if err != nil {
		return err
	}

	videos := make([]string, 0)

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".info.json") {
			continue
		}

		if f.Name() == playlistID+".info.json" {
			continue
		}

		videos = append(videos, strings.TrimSuffix(f.Name(), ".info.json"))
	}

	slog.Info("extracted playlist", "videos", videos)

	f, err := os.Open(filepath.Join(tempDir, playlistID+".info.json"))
	if err != nil {
		return err
	}

	var metadata playlistMetadata
	err = json.NewDecoder(f).Decode(&metadata)
	if err != nil {
		return err
	}

	description := fmt.Sprintf("%s, from playlist %s", metadata.Owner, playlistID)
	t, err := taskq.NewJsonTask(PriorityArchiveChannelInfo, TaskTypeArchiveChannelInfo, description, metadata.Owner)
	if err != nil {
		return err
	}

	err = taskq.Enqueue(t)
	if err != nil {
		return err
	}

	timestamp, err := time.ParseInLocation("20060102", metadata.Timestamp, time.Local)
	if err != nil {
		return err
	}

	for _, videoID := range videos {
		description := fmt.Sprintf("%s, from playlist %s", videoID, playlistID)
		task, err := taskq.NewJsonTask(PriorityArchiveVideo, TaskTypeArchiveVideo, description, videoID)
		if err != nil {
			return err
		}

		err = taskq.Enqueue(task)
		if err != nil {
			return err
		}
	}

	tx, err := db.DB().Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	err = db.Q().WithTx(tx).CreatePlaylist(
		context.Background(),
		db.CreatePlaylistParams{
			ID:          metadata.ID,
			Title:       metadata.Title,
			Description: metadata.Description,
			Timestamp:   timestamp,
			Owner:       metadata.Owner,
		})
	if err != nil {
		return err
	}

	for _, videoID := range videos {
		err = db.Q().WithTx(tx).CreatePlaylistVideo(
			context.Background(),
			db.CreatePlaylistVideoParams{
				Playlistid: playlistID,
				Videoid:    videoID,
			})
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
