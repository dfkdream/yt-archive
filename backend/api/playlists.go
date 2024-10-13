package api

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type playlistsHandler struct {
	DB *sql.DB
}

type Playlist struct {
	ID             string
	Title          string
	Description    string
	Timestamp      time.Time
	Owner          string
	OwnerThumbnail string
	ThumbnailVideo string
	Thumbnail      string
}

func (p playlistsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := `
	select p.id, p.title, p.description, p.timestamp, p.owner, c.thumbnail, v.id, v.thumbnail
	from playlists as p
	left join channels as c
	on p.owner=c.id
	left join playlist_video as pv
	on p.id=pv.playlistId
	left join videos as v
	on pv.videoId=v.id	
	where c.thumbnail not null
	and v.id not null
	order by p.rowid desc
	`

	rows, err := p.DB.Query(query)
	if err != nil {
		slog.Error("playlistsHandler error", "msg", err)
		writeError(w, http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	result := make([]Playlist, 0)
	var playlist Playlist
	var lastId string
	for rows.Next() {
		err = rows.Scan(&playlist.ID, &playlist.Title, &playlist.Description,
			&playlist.Timestamp, &playlist.Owner, &playlist.OwnerThumbnail, &playlist.ThumbnailVideo,
			&playlist.Thumbnail)

		if err != nil {
			slog.Error("playlistsHandler error", "msg", err)
			writeError(w, http.StatusInternalServerError)
			return
		}

		if lastId == playlist.ID {
			continue
		}

		result = append(result, playlist)
		lastId = playlist.ID
	}

	writeJson(w, result)
}

type playlistVideosHandler struct {
	DB *sql.DB
}

type VideoWithIndex struct {
	Video
	Index int
}

type PlaylistVideos struct {
	Playlist
	Videos []VideoWithIndex
}

func (p playlistVideosHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	query := `
	select p.id, p.title, p.description, p.timestamp, p.owner, c.thumbnail,
		v.id, v.title, v.description, v.timestamp, v.duration, v.owner, v.thumbnail, pv.sortIndex
	from playlists as p
	left join channels as c
	on p.owner=c.id
	left join playlist_video as pv
	on p.id=pv.playlistId
	left join videos as v
	on pv.videoId=v.id	
	where p.id=?
	order by pv.sortIndex asc, v.timestamp asc
	`

	rows, err := p.DB.Query(query, id)
	if err != nil {
		slog.Error("playlistVideosHandler error", "msg", err)
		writeError(w, http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var playlistVideos PlaylistVideos
	for rows.Next() {
		var video VideoWithIndex
		err = rows.Scan(
			&playlistVideos.ID, &playlistVideos.Title, &playlistVideos.Description,
			&playlistVideos.Timestamp, &playlistVideos.Owner, &playlistVideos.OwnerThumbnail,
			&video.ID, &video.Title, &video.Description, &video.Timestamp, &video.Duration,
			&video.Owner, &video.Thumbnail, &video.Index,
		)

		if err != nil {
			slog.Error("playlistVideos error", "msg", err)
			writeError(w, http.StatusInternalServerError)
			return
		}

		video.OwnerThumbnail = playlistVideos.OwnerThumbnail

		if playlistVideos.Thumbnail == "" {
			playlistVideos.ThumbnailVideo = video.ID
			playlistVideos.Thumbnail = video.Thumbnail
		}

		playlistVideos.Videos = append(playlistVideos.Videos, video)
	}

	writeJson(w, playlistVideos)
}

type playlistVideoIndexHandler struct {
	DB *sql.DB
}

func (p playlistVideoIndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	playlistID := mux.Vars(r)["pid"]
	videoID := mux.Vars(r)["vid"]
	var index int
	err := json.NewDecoder(r.Body).Decode(&index)
	if err != nil {
		slog.Error("playlistVideoIndexHandler error", "msg", err)
		writeError(w, http.StatusBadRequest)
		return
	}

	tx, err := p.DB.Begin()
	if err != nil {
		slog.Error("playlistVideoIndexHandler error", "msg", err)
		writeError(w, http.StatusInternalServerError)
		return
	}

	res, err := tx.Exec("update playlist_video set sortIndex=? where playlistId=? and videoId=?", index, playlistID, videoID)
	if err != nil {
		slog.Error("playlistVideoIndexHandler error", "msg", err)
		writeError(w, http.StatusInternalServerError)
		return
	}

	if r, _ := res.RowsAffected(); r != 1 {
		writeError(w, http.StatusNotFound)
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		slog.Error("playlistVideoIndexHandler error", "msg", err)
		writeError(w, http.StatusInternalServerError)
		return
	}

	writeJson(w, index)
}
