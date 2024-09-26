package api

import (
	"errors"
	"io/fs"
	"log/slog"
	"net/http"
	"path/filepath"
)

type StaticSiteServer struct {
	FS       http.FileSystem
	Fallback string
}

func (s *StaticSiteServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Clean(r.URL.Path)
	slog.Debug("StaticSiteServer request", "path", path)

	isdir, err := serveFile(w, r, s.FS, path)

	if errors.Is(err, fs.ErrNotExist) {
		s.serveFallback(w, r)
		return
	} else if err != nil {
		slog.Error("StaticSiteServer error", "msg", err)
		writeError(w, http.StatusInternalServerError)
		return
	}

	if isdir {
		s.serveIndex(w, r, path)
	}
}

func (s StaticSiteServer) serveIndex(w http.ResponseWriter, r *http.Request, path string) {
	_, err := serveFile(w, r, s.FS, filepath.Join(path, "index.html"))
	if errors.Is(err, fs.ErrNotExist) {
		s.serveFallback(w, r)
	} else if err != nil {
		slog.Error("serveIndex error", "msg", err)
		writeError(w, http.StatusInternalServerError)
	}
}

func (s StaticSiteServer) serveFallback(w http.ResponseWriter, r *http.Request) {
	_, err := serveFile(w, r, s.FS, s.Fallback)
	if err != nil {
		slog.Error("serveFallback error", "msg", err)
		writeError(w, http.StatusInternalServerError)
	}
}

func serveFile(w http.ResponseWriter, r *http.Request, FS http.FileSystem, path string) (bool, error) {
	f, err := FS.Open(path)
	if err != nil {
		return false, err
	}

	stat, err := f.Stat()
	if err != nil {
		return false, err
	}

	if stat.IsDir() {
		return true, nil
	}

	http.ServeContent(w, r, stat.Name(), stat.ModTime(), f)
	return false, nil
}

func writeError(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	w.Write([]byte(http.StatusText(code)))
}
