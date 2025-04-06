package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"yt-archive/db"

	"github.com/gorilla/mux"
)

type Video struct {
	db.Video
	Owner db.Channel
}

func videosHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Q().GetVideos(context.Background())
	if err != nil {
		slog.Error("videosHandler error", "msg", err)
		writeError(w, http.StatusInternalServerError)
		return
	}

	result := make([]Video, 0)
	for _, r := range rows {
		var video Video
		video.Video = r.Video
		video.Owner = r.Channel

		result = append(result, video)
	}

	writeJson(w, result)
}

func videoHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	row, err := db.Q().GetVideo(context.Background(), id)

	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, http.StatusNotFound)
		} else {
			slog.Error("videoHandler error", "msg", err)
			writeError(w, http.StatusInternalServerError)
		}
		return
	}

	var video Video
	video.Video = row.Video
	video.Owner = row.Channel

	writeJson(w, video)
}

func writeJson(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		w.Header().Del("Content-Type")
		slog.Error("writeJson error", "msg", err)
		writeError(w, http.StatusInternalServerError)
	}
}
