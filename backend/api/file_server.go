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
	slog.Info("fs request", "path", path)

	file, err := f.FS.Open(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			writeError(w, http.StatusNotFound)
			return
		}

		slog.Error("fs open", "error", err)
		return
	}

	stat, err := file.Stat()
	if err != nil {
		slog.Error("fs stat", "error", err)
		return
	}

	if stat.IsDir() {
		writeError(w, http.StatusNotFound)
		return
	}

	http.ServeContent(w, r, stat.Name(), stat.ModTime(), file)
}
