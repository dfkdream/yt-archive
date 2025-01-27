package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"yt-archive/taskq"

	"github.com/charmbracelet/huh"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	_, err := os.Stat("database/yt-archive.db")
	if errors.Is(err, os.ErrNotExist) {
		log.Fatal("Database file not found.")
	}

	db, err = sql.Open("sqlite3", "file:database/yt-archive.db?_journal_mode=WAL&_txlock=immediate")
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(1)

	q, err := taskq.New(db)
	if err != nil {
		log.Fatal(err)
	}

	taskq.SetDefaultQueue(q)
}

func main() {
	var selection int
	var functions []func() = []func(){
		showErroredTasks,
		enqueueAllErroredTasks,
		cancelAllErroredTasks,
		showFinishedTasks,
		deleteFinishedTasks,
		rebuildManifest,
		scanMissingVideoFiles,
	}

	fmt.Println("Welcome to yt-archive CLI!")

	for {
		err := huh.NewSelect[int]().
			Title("Main").
			Options(
				huh.NewOption("Show errored tasks", 0),
				huh.NewOption("Enqueue all errored tasks", 1),
				huh.NewOption("Cancel all errored tasks", 2),
				huh.NewOption("Show finished tasks", 3),
				huh.NewOption("Delete finished tasks", 4),
				huh.NewOption("Rebuild video manifest", 5),
				huh.NewOption("Scan missing video files", 6),
				huh.NewOption("Exit", -1),
			).
			Value(&selection).
			Run()

		if err != nil {
			log.Fatal(err)
		}

		if selection == -1 {
			fmt.Println("Bye!")
			os.Exit(0)
		}

		functions[selection]()
	}
}
