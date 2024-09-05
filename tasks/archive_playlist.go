package tasks

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
	"yt-archive/taskq"
)

const TASK_ARCHIVE_PLAYLIST = "ARCHIVE_PLAYLIST"

type ArchivePlaylistHandler struct {
	DB *sql.DB
}

func NewArchivePlaylistHandler(db *sql.DB) (ArchivePlaylistHandler, error) {
	_, err := db.Exec("create table if not exists playlists (id text primary key unique, title text, description text, timestamp timestamp, owner text)")
	if err != nil {
		return ArchivePlaylistHandler{}, err
	}

	_, err = db.Exec("create table if not exists playlist_video (playlistId text, videoId text, unique (playlistId, videoId), primary key (playlistId, videoId))")
	if err != nil {
		return ArchivePlaylistHandler{}, err
	}

	return ArchivePlaylistHandler{DB: db}, nil
}

type playlistMetadata struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Owner       string `json:"uploader_id"`
	Timestamp   string `json:"modified_date"`
}

func (a ArchivePlaylistHandler) Handler(task *taskq.Task) error {
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

	timestamp, err := time.ParseInLocation("20060102", metadata.Timestamp, time.Local)
	if err != nil {
		return err
	}

	for _, videoID := range videos {
		task, err := taskq.NewJsonTask(PRIORITY_ARCHIVE_VIDEO, TASK_ARCHIVE_VIDEO, playlistID+" - "+videoID, videoID)
		if err != nil {
			return err
		}

		err = taskq.Enqueue(task)
		if err != nil {
			return err
		}
	}

	tx, err := a.DB.Begin()
	if err != nil {
		return err
	}

	query := `
	insert into playlists (id, title, description, timestamp, owner) values (?, ?, ?, ?, ?)
	on conflict(id) do update set timestamp=excluded.timestamp
	`

	_, err = tx.Exec(query, metadata.ID, metadata.Title, metadata.Description, timestamp, metadata.Owner)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, videoID := range videos {
		_, err = tx.Exec("insert into playlist_video (playlistId, videoId) values (?, ?) on conflict(playlistId, videoId) do nothing",
			playlistID, videoID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
