package api

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type channelsHandler struct {
	DB *sql.DB
}

type Channel struct {
	ID          string
	Title       string
	Description string
	Thumbnail   string
}

func (c channelsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rows, err := c.DB.Query("select id, title, description, thumbnail from channels")
	if err != nil {
		slog.Error("channelsHandler error", "msg", err)
		writeError(w, http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	result := make([]Channel, 0)
	var channel Channel
	for rows.Next() {
		err = rows.Scan(&channel.ID, &channel.Title, &channel.Description, &channel.Thumbnail)
		if err != nil {
			slog.Error("channelsHandler error", "msg", err)
			writeError(w, http.StatusInternalServerError)
			return
		}

		result = append(result, channel)
	}

	writeJson(w, result)
}

type channelVideosHandler struct {
	DB *sql.DB
}

type ChannelVideos struct {
	Channel
	Videos []Video
}

func (c channelVideosHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	query := `
	select c.id, c.title, c.description, c.thumbnail,
		v.id, v.title, v.description, timestamp, duration, owner, v.thumbnail
	from channels as c
	left join videos as v
	on v.owner = c.id
	where c.id=?
	order by v.rowid desc
	`

	rows, err := c.DB.Query(query, id)
	if err != nil {
		slog.Error("channelVideosHandler error", "msg", err)
		writeError(w, http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var channelVideos ChannelVideos
	for rows.Next() {
		var video Video
		err = rows.Scan(
			&channelVideos.ID, &channelVideos.Title, &channelVideos.Description, &channelVideos.Thumbnail,
			&video.ID, &video.Title, &video.Description, &video.Timestamp, &video.Duration, &video.Owner, &video.Thumbnail,
		)

		if err != nil {
			slog.Error("channelVideos error", "msg", err)
			writeError(w, http.StatusInternalServerError)
			return
		}

		video.OwnerThumbnail = channelVideos.ID

		channelVideos.Videos = append(channelVideos.Videos, video)
	}

	if channelVideos.Videos == nil {
		writeError(w, http.StatusNotFound)
		return
	}

	writeJson(w, channelVideos)
}
