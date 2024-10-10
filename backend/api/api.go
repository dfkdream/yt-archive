package api

import (
	"database/sql"
	"io/fs"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func New(db *sql.DB, distFS fs.FS) http.Handler {
	r := mux.NewRouter()

	r.Path("/api/videos").
		Methods(http.MethodGet).
		Handler(videosHandler{DB: db})

	r.Path("/api/videos/{id}").
		Methods(http.MethodGet).
		Handler(videoHandler{DB: db})

	r.PathPrefix("/api/videos/").
		Methods(http.MethodGet).
		Handler(http.StripPrefix("/api/videos/", &FileServer{
			FS: http.FS(os.DirFS("videos")),
		}))

	r.Path("/api/channels").
		Methods(http.MethodGet).
		Handler(channelsHandler{DB: db})

	r.Path("/api/channels/{id}").
		Methods(http.MethodGet).
		Handler(channelVideosHandler{DB: db})

	r.PathPrefix("/api/channels/").
		Methods(http.MethodGet).
		Handler(http.StripPrefix("/api/channels/", &FileServer{
			FS: http.FS(os.DirFS("channels")),
		}))

	r.Path("/api/playlists").
		Methods(http.MethodGet).
		Handler(playlistsHandler{DB: db})

	r.Path("/api/playlists/{id}").
		Methods(http.MethodGet).
		Handler(playlistVideosHandler{DB: db})

	r.Path("/api/playlists/{pid}/video/{vid}/index").
		Methods(http.MethodPost).
		Headers("Content-Type", "application/json").
		Handler(playlistVideoIndexHandler{DB: db})

	r.Path("/api/tasks").
		Methods(http.MethodGet).
		Handler(tasksHandler{DB: db})

	r.Path("/api/tasks").
		Methods(http.MethodPost).
		Headers("Content-Type", "application/json").
		HandlerFunc(enqueTaskHandler)

	r.Path("/api/tasks/{id}").
		Methods(http.MethodGet).
		Handler(taskHandler{DB: db})

	r.PathPrefix("/").
		Methods(http.MethodGet).
		Handler(&StaticSiteServer{
			FS:       http.FS(distFS),
			Fallback: "fallback.html",
		})

	return r
}
