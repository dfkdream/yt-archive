package tasks

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"yt-archive/taskq"
)

const TaskTypeArchiveChannelInfo = "ARCHIVE_CHANNEL_INFO"

type ArchiveChannelInfoHandler struct {
	DB *sql.DB
}

func NewArchiveChannelInfoHandler(db *sql.DB) (ArchiveChannelInfoHandler, error) {
	_, err := db.Exec("create table if not exists channels (id text primary key, title text, description text, thumbnail text)")
	if err != nil {
		return ArchiveChannelInfoHandler{}, err
	}

	return ArchiveChannelInfoHandler{DB: db}, nil
}

type channelMetadata struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (a ArchiveChannelInfoHandler) Handler(task *taskq.Task) error {
	var channelID string
	err := json.Unmarshal(task.Payload, &channelID)
	if err != nil {
		return err
	}

	r := a.DB.QueryRow("select count(id) from channels where id=?", channelID)
	var n int
	err = r.Scan(&n)
	if err != nil {
		return err
	}

	if n > 0 {
		slog.Info("skipping archived channel", "id", channelID)
		return nil
	}

	tempDir, err := os.MkdirTemp("", channelID+"_*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	channelUrl, err := url.JoinPath("https://www.youtube.com", channelID)
	if err != nil {
		return err
	}

	err = Exec("yt-dlp", "--write-info-json", "--skip-download", "--write-thumbnail", "--playlist-items", "0", "-o", "%(id)s.%(ext)s", "--paths", tempDir, channelUrl)
	if err != nil {
		return err
	}

	thumbnail, err := copyThumbnail(tempDir)
	if err != nil {
		return err
	}

	slog.Info("downloaded thumbnail", "filename", thumbnail)

	f, err := os.Open(filepath.Join(tempDir, channelID+".info.json"))
	if err != nil {
		return err
	}
	defer f.Close()

	var metadata channelMetadata
	err = json.NewDecoder(f).Decode(&metadata)
	if err != nil {
		return err
	}

	_, err = a.DB.Exec("insert into channels (id, title, description, thumbnail) values (?, ?, ?, ?)",
		metadata.ID, metadata.Title, metadata.Description, thumbnail)

	if err != nil {
		return err
	}

	return nil
}
