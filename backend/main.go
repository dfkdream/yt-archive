package main

import (
	"database/sql"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"yt-archive/api"
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

func main() {
	err := os.MkdirAll("database", os.FileMode(0o700))
	if err != nil {
		log.Fatal(err)
	}

	migrationRequired := true
	_, err = os.Stat("database/yt-archive.db")
	if errors.Is(err, os.ErrNotExist) {
		migrationRequired = false
	}

	db, err := sql.Open("sqlite3", "file:database/yt-archive.db?_journal_mode=WAL&_txlock=immediate")
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(1)

	migrate(db, migrationRequired)

	q, err := taskq.New(db)
	if err != nil {
		log.Fatal(err)
	}

	taskq.SetDefaultQueue(q)

	archiveVideo, err := tasks.NewArchiveVideoHandler(db)
	if err != nil {
		log.Fatal(err)
	}

	archivePlaylist, err := tasks.NewArchivePlaylistHandler(db)
	if err != nil {
		log.Fatal(err)
	}

	archiveChannelInfo, err := tasks.NewArchiveChannelInfoHandler(db)
	if err != nil {
		log.Fatal(err)
	}

	taskq.Handler(tasks.TaskTypeArchiveVideo, archiveVideo.Handler)
	taskq.Handler(tasks.TaskTypeDownloadMedia, tasks.DownloadMediaHandler)
	taskq.Handler(tasks.TaskTypeArchivePlaylist, archivePlaylist.Handler)
	taskq.Handler(tasks.TaskTypeArchiveChannelInfo, archiveChannelInfo.Handler)

	go taskq.Start()

	http.Handle("/", api.New(db))

	addr := os.Getenv("YT_ARCHIVE_ADDR")
	if addr == "" {
		addr = "0.0.0.0:80"
	}

	log.Printf("YT_ARCHIVE_ADDR: %s\n", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}
