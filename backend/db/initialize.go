package db

import (
	"database/sql"
	_ "embed"
	"os"
)

//go:embed sql/schema.sql
var schema string

var db *sql.DB
var queries *Queries

func DB() *sql.DB {
	if db == nil {
		err := os.MkdirAll("database", os.FileMode(0o700))
		if err != nil {
			panic(err)
		}

		db, err = sql.Open("sqlite3", "file:database/yt-archive.db?_journal_mode=WAL&_txlock=immediate")
		if err != nil {
			panic(err)
		}

		db.SetMaxOpenConns(1)

		_, err = db.Exec(schema)
		if err != nil {
			panic(err)
		}
	}

	return db
}

func Q() *Queries {
	if queries == nil {
		queries = New(DB())
	}

	return queries
}
