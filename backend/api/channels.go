package api

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"yt-archive/db"

	"github.com/gorilla/mux"
)

func channelsHandler(w http.ResponseWriter, r *http.Request) {
	channels, err := db.Q().GetChannels(context.Background())
	if err != nil {
		slog.Error("channelsHandler error", "msg", err)
		writeError(w, http.StatusInternalServerError)
		return
	}

	writeJson(w, channels)
}

type ChannelVideos struct {
	db.Channel
	Videos []Video
}

func channelVideosHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var channelVideos ChannelVideos
	channel, err := db.Q().GetChannel(context.Background(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, http.StatusNotFound)
			return
		}

		slog.Error("channelVideosHandler error", "msg", err)
		writeError(w, http.StatusInternalServerError)
		return
	}
	channelVideos.Channel = channel

	rows, err := db.Q().GetChannelVideos(context.Background(), id)
	if err != nil {
		slog.Error("channelVideosHandler error", "msg", err)
		writeError(w, http.StatusInternalServerError)
		return
	}

	for _, r := range rows {
		var video Video
		video.Video = r.Video
		video.Owner = r.Channel
		channelVideos.Videos = append(channelVideos.Videos, video)
	}

	if channelVideos.Videos == nil {
		channelVideos.Videos = []Video{}
	}

	writeJson(w, channelVideos)
}
