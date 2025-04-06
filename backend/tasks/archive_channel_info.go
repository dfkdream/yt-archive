package tasks

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"yt-archive/db"
	"yt-archive/taskq"
)

const TaskTypeArchiveChannelInfo = "ARCHIVE_CHANNEL_INFO"

type channelMetadata struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func ArchiveChannelInfoHandler(task *taskq.Task) error {
	var channelID string
	err := json.Unmarshal(task.Payload, &channelID)
	if err != nil {
		return err
	}

	n, err := db.Q().GetChannelCount(context.Background(), channelID)
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

	err = db.Q().CreateChannel(
		context.Background(),
		db.CreateChannelParams{
			ID:          metadata.ID,
			Title:       metadata.Title,
			Description: metadata.Description,
			Thumbnail:   thumbnail,
		})

	if err != nil {
		return err
	}

	return nil
}
