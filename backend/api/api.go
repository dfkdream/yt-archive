package api

import (
	"net/http"
	"os"
)

func New() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", &StaticSiteServer{
		FS:       http.FS(os.DirFS("dist")),
		Fallback: "fallback.html",
	})

	mux.Handle("/api/videos/", http.StripPrefix("/api/videos", &FileServer{
		FS: http.FS(os.DirFS("videos")),
	}))

	mux.Handle("/api/channels/", http.StripPrefix("/api/channels", &FileServer{
		FS: http.FS(os.DirFS("channels")),
	}))

	return mux
}
