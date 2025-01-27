package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"yt-archive/mpd"
	"yt-archive/taskq"
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

func enqueueAllErroredTasks() {
	rowsAffected, err := execRowsAffected("update tasks set status=0 where status=4")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d tasks updated.\n", rowsAffected)
}

func cancelAllErroredTasks() {
	rowsAffected, err := execRowsAffected("update tasks set status=2 where status=4")
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

var videoIDRegex = regexp.MustCompile("^[A-Za-z0-9_-]{11}$")

func videoIDValidator(s string) error {
	if !videoIDRegex.MatchString(s) {
		return fmt.Errorf("invalid video ID")
	}

	return nil
}

func rebuildManifest() {
	videoID := ""

	err := huh.NewInput().
		Title("Rebuild manifest").
		Description("Enter video ID").
		Validate(videoIDValidator).
		Value(&videoID).
		Run()

	if err != nil {
		log.Fatal(err)
	}

	videoPath := filepath.Join("videos", videoID)
	if !isDirExist(videoPath) {
		log.Printf("Video %s not found\n", videoID)
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

func scanMissingVideoFiles() {
	videoID := ""

	err := huh.NewInput().
		Title("Scan missing video files").
		Description("Enter video ID").
		Validate(videoIDValidator).
		Value(&videoID).
		Run()

	if err != nil {
		log.Fatal(err)
	}

	videoPath := filepath.Join("videos", videoID)
	if !isDirExist(videoPath) {
		log.Printf("Video %s not found\n", videoID)
		return
	}

	metadata, _, err := tasks.DownloadVideoMetadata(videoID, false)
	if err != nil {
		log.Println(err)
		return
	}

	tasksToEnqueue := make([]*taskq.Task, 0)

	payload := tasks.DownloadMediaPayload{
		VideoID:      videoID,
		Format:       "bestaudio",
		OutputPath:   filepath.Join(videoPath, videoID+tasks.AUDIO_FILE_SUFFIX),
		SkipEncoding: false,
	}

	if !isFileExist(payload.OutputPath) {
		// Enqueue audio download task
		t, err := taskq.NewJsonTask(
			tasks.PriorityDownloadAudio,
			tasks.TaskTypeDownloadMedia,
			videoID+", bestaudio (CLI)",
			payload,
		)

		if err != nil {
			log.Println(err)
			return
		}

		tasksToEnqueue = append(tasksToEnqueue, t)
	}

	formats := tasks.SelectVideoFormats(metadata.Formats)
	for _, v := range formats {
		videoFilePath := filepath.Join(videoPath, fmt.Sprintf("%s_%d"+tasks.MEDIA_FILE_SUFFIX, videoID, v.Height))
		if !isFileExist(videoFilePath) {
			// Enqueue video download task
			payload.Format = v.FormatID
			payload.OutputPath = videoFilePath
			payload.SkipEncoding = tasks.CanSkipEncoding(v)

			description := fmt.Sprintf("%s, %d, %s", videoID, v.Height, v.VideoCodec)
			if !payload.SkipEncoding {
				description += " (Encoding required)"
			}
			description += " (CLI)"

			t, err := taskq.NewJsonTask(
				tasks.CalculateVideoPriority(v),
				tasks.TaskTypeDownloadMedia,
				description,
				payload,
			)

			if err != nil {
				log.Println(err)
				return
			}

			tasksToEnqueue = append(tasksToEnqueue, t)
		}
	}

	if len(tasksToEnqueue) == 0 {
		fmt.Println("Everything is OK")
		return
	}

	fmt.Println("Tasks to be enqueued:")
	for _, t := range tasksToEnqueue {
		fmt.Println(t.ID, t.Description)
	}

	enqueueNow := true
	huh.NewConfirm().
		Title("Would you like to enqueue these tasks?").
		Value(&enqueueNow).
		Run()

	if !enqueueNow {
		return
	}

	for _, t := range tasksToEnqueue {
		err = taskq.Enqueue(t)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func isDirExist(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}

	return stat.IsDir()
}

func isFileExist(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !stat.IsDir()
}
