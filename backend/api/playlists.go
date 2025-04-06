package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"yt-archive/db"

	"github.com/gorilla/mux"
)

type Playlist struct {
	db.Playlist
	Owner          db.Channel
	ThumbnailVideo db.Video
}

func playlistsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Q().GetPlaylists(context.Background())
	if err != nil {
		slog.Error("playlistsHandler error", "msg", err)
		writeError(w, http.StatusInternalServerError)
		return
	}

	result := make([]Playlist, 0)
	var playlist Playlist
	var lastId string

	for _, r := range rows {
		if lastId == r.Playlist.ID {
			continue
		}
		lastId = r.Playlist.ID

		playlist.Playlist = r.Playlist
		playlist.Owner = r.Channel
		playlist.ThumbnailVideo = r.Video

		result = append(result, playlist)
	}

	writeJson(w, result)
}

type VideoWithIndex struct {
	Video
	Index int
}

type PlaylistVideos struct {
	db.Playlist
	Owner  db.Channel
	Videos []VideoWithIndex
}

func playlistVideosHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	rows, err := db.Q().GetPlaylistVideos(context.Background(), id)
	if err != nil {
		slog.Error("playlistVideosHandler error", "msg", err)
		writeError(w, http.StatusInternalServerError)
		return
	}

	var playlistVideos PlaylistVideos
	for _, r := range rows {
		playlistVideos.Playlist = r.Playlist
		playlistVideos.Owner = r.Channel

		var video VideoWithIndex
		video.Video.Video = r.Video
		video.Owner = r.Channel_2
		video.Index = int(r.Sortindex.Int64)

		playlistVideos.Videos = append(playlistVideos.Videos, video)
	}

	if playlistVideos.Videos == nil {
		writeError(w, http.StatusNotFound)
		return
	}

	writeJson(w, playlistVideos)
}

func playlistVideoIndexHandler(w http.ResponseWriter, r *http.Request) {
	playlistID := mux.Vars(r)["pid"]
	videoID := mux.Vars(r)["vid"]
	var index int64
	err := json.NewDecoder(r.Body).Decode(&index)
	if err != nil {
		slog.Error("playlistVideoIndexHandler error", "msg", err)
		writeError(w, http.StatusBadRequest)
		return
	}

	tx, err := db.DB().Begin()
	if err != nil {
		slog.Error("playlistVideoIndexHandler error", "msg", err)
		writeError(w, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	rowsAffected, err := db.Q().WithTx(tx).UpdatePlaylistVideoIndex(
		context.Background(),
		db.UpdatePlaylistVideoIndexParams{
			Sortindex:  index,
			Playlistid: playlistID,
			Videoid:    videoID,
		})

	if err != nil {
		slog.Error("playlistVideoIndexHandler error", "msg", err)
		writeError(w, http.StatusInternalServerError)
		return
	}

	if rowsAffected != 1 {
		writeError(w, http.StatusNotFound)
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
