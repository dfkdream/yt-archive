package api

import (
	"errors"
	"io/fs"
	"log/slog"
	"net/http"
	"path/filepath"
)

type FileServer struct {
	FS http.FileSystem
}

func (f *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Clean(r.URL.Path)
	slog.Info("FileServer request", "path", path)

	isDir, err := serveFile(w, r, f.FS, path)

	if isDir || errors.Is(err, fs.ErrNotExist) {
		writeError(w, http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error("FileServer error", "msg", err)
		writeError(w, http.StatusInternalServerError)
	}
}
