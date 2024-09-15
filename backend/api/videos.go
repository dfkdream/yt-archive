package api

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type videosHandler struct {
	DB *sql.DB
}

type Video struct {
	ID             string
	Title          string
	Description    string
	Timestamp      time.Time
	Duration       string
	Owner          string
	Thumbnail      string
	OwnerThumbnail string
}

func (v videosHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := `
	select videos.id, videos.title, videos.description, timestamp, duration, owner, videos.thumbnail, channels.thumbnail
	from videos
	left join channels
	on videos.owner = channels.id
	order by videos.rowid desc
	`

	rows, err := v.DB.Query(query)
	if err != nil {
		slog.Error("videosHandler error", "msg", err)
		writeError(w, http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	result := make([]Video, 0)
	var video Video
	for rows.Next() {
		err = rows.Scan(&video.ID, &video.Title, &video.Description,
			&video.Timestamp, &video.Duration, &video.Owner, &video.Thumbnail, &video.OwnerThumbnail)

		if err != nil {
			slog.Error("videosHandler error", "msg", err)
			writeError(w, http.StatusInternalServerError)
			return
		}

		result = append(result, video)
	}

	writeJson(w, result)
}

type videoHandler struct {
	DB *sql.DB
}

func (v videoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	query := `
	select videos.id, videos.title, videos.description, timestamp, duration, owner, videos.thumbnail, channels.thumbnail
	from videos
	left join channels
	on videos.owner = channels.id
	where videos.id=?
	`

	row := v.DB.QueryRow(query, id)

	var video Video
	err := row.Scan(&video.ID, &video.Title, &video.Description,
		&video.Timestamp, &video.Duration, &video.Owner, &video.Thumbnail, &video.OwnerThumbnail)

	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, http.StatusNotFound)
		} else {
			slog.Error("videoHandler error", "msg", err)
			writeError(w, http.StatusInternalServerError)
		}
		return
	}

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
