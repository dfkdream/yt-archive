package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
)

var queries = []string{
	"alter table playlist_video add column sortIndex integer",
	"update playlist_video set sortIndex=0 where sortIndex is null",
}

func migrate(DB *sql.DB, required bool) {
	var currentVersion int
	r := DB.QueryRow("pragma user_version")
	err := r.Scan(&currentVersion)
	if err != nil {
		log.Fatal(err)
	}

	if !required {
		log.Printf("New database detected. Skipping migrations.")
	}

	if required {
		log.Printf("Current database version is %d.\n", currentVersion)

		if currentVersion == len(queries) {
			log.Println("Database is up to date.")
			return
		}

		if currentVersion > len(queries) {
			log.Printf("Current version is higher than defined, which is %d.\n", len(queries))
			log.Println("Your database might be corrupted. Skipping migration.")
			return
		}

		for _, q := range queries[currentVersion:] {
			slog.Info("Running migration", "query", q)
			_, err = DB.Exec(q)
			if err != nil {
				log.Fatal("migration error: ", err)
			}
		}
	}

	_, err = DB.Exec(fmt.Sprintf("pragma user_version=%d", len(queries)))
	if err != nil {
		log.Fatal("migration error: ", err)
	}
}
