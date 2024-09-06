package api

import (
	"errors"
	"io/fs"
	"log/slog"
	"net/http"
	"path/filepath"
)

type StaticSiteServer struct {
	FS http.FileSystem
}

func (s *StaticSiteServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Clean(r.URL.Path)
	slog.Info("sss request", "path", path)

	f, err := s.FS.Open(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			f, err = findIndex(s.FS, path)
			if err != nil {
				slog.Error("StaticSiteServer Open findIndex", "error", err)
				writeError(w, http.StatusInternalServerError)
				return
			}
		} else {
			slog.Error("StaticSiteServer Open", "error", err)
			writeError(w, http.StatusInternalServerError)
			return
		}
	}

	stat, err := f.Stat()
	if err != nil {
		slog.Error("StaticSiteServer Stat", "error", err)
		writeError(w, http.StatusInternalServerError)
		return
	}

	if stat.IsDir() {
		f.Close()
		f, err = findIndex(s.FS, path)
		if err != nil {
			slog.Error("StaticSiteServer findIndex", "error", err)
			writeError(w, http.StatusInternalServerError)
			return
		}
	}

	http.ServeContent(w, r, stat.Name(), stat.ModTime(), f)
}

func writeError(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	w.Write([]byte(http.StatusText(code)))
}

const indexFile = "index.html"

func findIndex(targetFs http.FileSystem, path string) (http.File, error) {
	path = filepath.Clean(path)

	slog.Info("findIndex", "path", path)

	f, err := targetFs.Open(filepath.Join(path, indexFile))
	if err == nil {
		stat, err := f.Stat()
		if err != nil {
			return nil, err
		}

		if !stat.IsDir() {
			return f, nil
		}
	}

	if !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}

	if path == "/" {
		return nil, fs.ErrNotExist
	}

	return findIndex(targetFs, filepath.Join(path, ".."))
}
