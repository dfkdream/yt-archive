package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"yt-archive/mpd"
	"yt-archive/tasks"

	"github.com/charmbracelet/huh"
	"github.com/google/uuid"
)

func showErroredTasks() {
	rows, err := db.Query("select id, description from tasks where status=4")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	fmt.Println("Errored tasks:")
	for rows.Next() {
		var id uuid.UUID
		var description string
		err = rows.Scan(&id, &description)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(id, description)
	}
}

func resetAllErroredTasks() {
	rowsAffected, err := execRowsAffected("update tasks set status=0 where status=4")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d tasks updated.\n", rowsAffected)
}

func showFinishedTasks() {
	rows, err := db.Query("select id, description from tasks where status=3")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	fmt.Println("Finished tasks:")

	count := 0
	const limit = 10

	for rows.Next() {
		var id uuid.UUID
		var description string
		err = rows.Scan(&id, &description)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(id, description)

		count++
		if count > limit {
			fmt.Println("... and more. Stopping here.")
			return
		}
	}
}

func positiveIntValidator(s string) error {
	n, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("enter positive integer")
	}

	if n < 0 {
		return fmt.Errorf("enter positive integer")
	}

	return nil
}

func deleteFinishedTasks() {
	strPreserveN := "10"

	err := huh.NewInput().
		Title("Delete finished tasks").
		Description("How many finished tasks do you want to preserve?").
		Validate(positiveIntValidator).
		Value(&strPreserveN).
		Run()

	if err != nil {
		log.Fatal(err)
	}

	preserveN, err := strconv.Atoi(strPreserveN)
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err := execRowsAffected("delete from tasks where status=3 and id not in (select id from tasks where status=3 order by id desc limit ?)", preserveN)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d tasks deleted.\n", rowsAffected)
}

func rebuildManifest() {
	videoID := ""

	err := huh.NewInput().
		Title("Rebuild manifest").
		Description("Enter video ID").
		Value(&videoID).
		Run()

	if err != nil {
		log.Fatal(err)
	}

	videoPath := filepath.Join("videos", videoID)
	if _, err := os.Stat(videoPath); err != nil {
		log.Println(err)
		return
	}

	tempDir, err := os.MkdirTemp("", videoID+"_*")
	if err != nil {
		log.Println(err)
		return
	}
	defer os.RemoveAll(tempDir)

	files, err := os.ReadDir(videoPath)
	if err != nil {
		log.Println(err)
		return
	}

	videoManifest := filepath.Join(tempDir, "video.mpd")
	masterManifest := filepath.Join(tempDir, "master.mpd")

	for _, f := range files {
		if strings.HasSuffix(f.Name(), tasks.MEDIA_FILE_SUFFIX) {
			tasks.BuildManifest(
				filepath.Join(videoPath, f.Name()),
				videoManifest,
			)

			if _, err := os.Stat(masterManifest); err != nil {
				if errors.Is(err, fs.ErrNotExist) {
					err = os.Rename(videoManifest, masterManifest)
					if err != nil {
						log.Println(err)
						return
					}
				} else {
					log.Println(err)
					return
				}
			} else {
				master, err := mpd.FromFile(masterManifest)
				if err != nil {
					log.Println(err)
					return
				}

				video, err := mpd.FromFile(videoManifest)
				if err != nil {
					log.Println(err)
					return
				}

				mpd.Merge(master, video).WriteFile(masterManifest)
			}
		}
	}

	finalManifest := filepath.Join(videoPath, videoID+".mpd")
	err = os.Rename(masterManifest, finalManifest)
	if err != nil {
		log.Println(err)
	}

	log.Println("Done! Manifest written to", finalManifest)
}
