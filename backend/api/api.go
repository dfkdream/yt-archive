package api

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func New(db *sql.DB) http.Handler {
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

	r.PathPrefix("/").
		Methods(http.MethodGet).
		Handler(&StaticSiteServer{
			FS:       http.FS(os.DirFS("dist")),
			Fallback: "fallback.html",
		})

	return r
}
