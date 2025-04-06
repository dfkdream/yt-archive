package main

import (
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"
	"yt-archive/api"
	"yt-archive/db"
	"yt-archive/taskq"
	"yt-archive/tasks"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	logLevel := os.Getenv("YT_ARCHIVE_LOG_LEVEL")

	switch logLevel {
	case "ERROR":
		slog.SetLogLoggerLevel(slog.LevelError)
	case "INFO":
		slog.SetLogLoggerLevel(slog.LevelInfo)
	case "DEBUG":
		slog.SetLogLoggerLevel(slog.LevelDebug)
	default:
		slog.SetLogLoggerLevel(slog.LevelInfo)
		logLevel = "INFO"
	}

	log.Printf("YT_ARCHIVE_LOG_LEVEL: %s\n", logLevel)
}

func entrypoint(distFS fs.FS) {
	q, err := taskq.New(db.DB())
	if err != nil {
		log.Fatal(err)
	}

	err = q.ResetRunningTasks()
	if err != nil {
		log.Fatal(err)
	}

	taskq.SetDefaultQueue(q)

	taskq.Handler(tasks.TaskTypeArchiveVideo, tasks.ArchiveVideoHandler)
	taskq.Handler(tasks.TaskTypeDownloadMedia, tasks.DownloadMediaHandler)
	taskq.Handler(tasks.TaskTypeArchivePlaylist, tasks.ArchivePlaylistHandler)
	taskq.Handler(tasks.TaskTypeArchiveChannelInfo, tasks.ArchiveChannelInfoHandler)

	go taskq.Start()

	http.Handle("/", api.New(db.DB(), distFS))

	addr := os.Getenv("YT_ARCHIVE_ADDR")
	if addr == "" {
		addr = "0.0.0.0:80"
	}

	log.Printf("YT_ARCHIVE_ADDR: %s\n", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}
