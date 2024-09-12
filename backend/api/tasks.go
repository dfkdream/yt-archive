package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"yt-archive/taskq"
	"yt-archive/tasks"

	"github.com/gorilla/mux"
)

type tasksHandler struct {
	DB *sql.DB
}

func (t tasksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rows, err := t.DB.Query("select id, status, priority, type, description, payload from tasks order by status asc, priority desc, id asc")
	if err != nil {
		slog.Error("tasksHandler error", "msg", err)
		writeError(w, http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	result := make([]taskq.Task, 0)
	var task taskq.Task
	for rows.Next() {
		err = rows.Scan(&task.ID, &task.Status, &task.Priority,
			&task.Type, &task.Description, &task.Payload)

		if err != nil {
			slog.Error("taskHandler error", "msg", err)
			writeError(w, http.StatusInternalServerError)
			return
		}

		result = append(result, task)
	}

	writeJson(w, result)
}

type taskHandler struct {
	DB *sql.DB
}

func (t taskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	row := t.DB.QueryRow("select id, status, priority, type, description, payload from tasks where id=?", id)

	var task taskq.Task
	err := row.Scan(&task.ID, &task.Status, &task.Priority,
		&task.Type, &task.Description, &task.Payload)

	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, http.StatusNotFound)
		} else {
			slog.Error("taskHandler error", "msg", err)
			writeError(w, http.StatusInternalServerError)
		}

		return
	}

	writeJson(w, task)
}

const (
	typeVideo = iota
	typePlaylist
)

type enqueueTaskParams struct {
	Type int
	ID   string
}

var videoIDRegex = regexp.MustCompile("^[A-Za-z0-9_-]{11}$")
var playlistIDRegex = regexp.MustCompile("^PL[A-Za-z0-9_-]{32}$")

func enqueTaskHandler(w http.ResponseWriter, r *http.Request) {
	var params enqueueTaskParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		slog.Error("enqueueTaskHandler error", "msg", err)
		writeError(w, http.StatusBadRequest)
		return
	}

	var tasktype string
	switch params.Type {
	case typeVideo:
		tasktype = tasks.TaskTypeArchiveVideo
	case typePlaylist:
		tasktype = tasks.TaskTypeArchivePlaylist
	default:
		writeError(w, http.StatusBadRequest)
		return
	}

	if params.Type == typeVideo && !videoIDRegex.MatchString(params.ID) {
		writeError(w, http.StatusBadRequest)
		return
	}

	if params.Type == typePlaylist && !playlistIDRegex.MatchString(params.ID) {
		writeError(w, http.StatusBadRequest)
		return
	}

	task, err := taskq.NewJsonTask(tasks.PriorityHighest, tasktype,
		fmt.Sprintf("Request from %s", r.RemoteAddr), params.ID)
	if err != nil {
		slog.Error("enqueueTaskHandler error", "msg", err)
		writeError(w, http.StatusInternalServerError)
		return
	}

	err = taskq.Enqueue(task)
	if err != nil {
		slog.Error("enqueueTaskHandler error", "msg", err)
		writeError(w, http.StatusInternalServerError)
		return
	}

	writeJson(w, task.ID)
}
