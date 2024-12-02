package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

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
}

func main() {
	var selection int
	var functions []func() = []func(){
		showErroredTasks,
		resetAllErroredTasks,
		showFinishedTasks,
		deleteFinishedTasks,
	}

	fmt.Println("Welcome to yt-archive CLI!")

	for {
		err := huh.NewSelect[int]().
			Title("Main").
			Options(
				huh.NewOption("Show errored tasks", 0),
				huh.NewOption("Reset all errored tasks", 1),
				huh.NewOption("Show finished tasks", 2),
				huh.NewOption("Delete finished tasks", 3),
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
