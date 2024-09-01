package taskq

import (
	"database/sql"
	"errors"
	"log"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

const (
	statusQueued = iota
	statusRunning
	statusCancelled
	statusFinished
	statusError
)

type Task struct {
	ID          uuid.UUID
	Status      int
	Priority    int
	Type        string
	Description string
	Payload     string
}

type TaskHandler func(task *Task) error

var db *sql.DB

var handlers map[string]TaskHandler

var fallbackHandler TaskHandler = func(task *Task) error {
	slog.Info("ignoring task with no matching handler", "type", task.Type)
	return nil
}

var cond = sync.NewCond(new(sync.Mutex))

func init() {
	var err error
	db, err = sql.Open("sqlite3", "file:taskq.db")
	db.SetMaxOpenConns(1)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("opened sqlite database", "filename", "taskq.db")

	_, err = db.Exec("create table if not exists tasks (id text primary key, status integer, priority integer, type text, description text, payload text) strict")
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("created table tasks")

	handlers = make(map[string]TaskHandler)
}

func NewTask(priority int, tasktype, description string, payload string) (*Task, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return &Task{
		ID:          uuid,
		Status:      statusQueued,
		Priority:    priority,
		Type:        tasktype,
		Description: description,
		Payload:     payload,
	}, nil
}

func Enqueue(task *Task) error {

	_, err := db.Exec("insert into tasks (id, status, priority, type, description, payload) values (?, ?, ?, ?, ?, ?)",
		task.ID, task.Status, task.Priority, task.Type, task.Description, task.Payload)

	if err == nil {
		cond.L.Lock()
		slog.Info("task enqueued", "type", task.Type, "description", task.Description)
		cond.Signal()
		cond.L.Unlock()
	}

	return err
}

func Dispatch() {
	row := db.QueryRow("select id, status, priority, type, description, payload from tasks where status=? order by priority desc, id asc limit 1", statusQueued)

	var task Task
	err := row.Scan(&task.ID, &task.Status, &task.Priority, &task.Type, &task.Description, &task.Payload)
	if err != nil {
		if err == sql.ErrNoRows {
			//wait until enqueue
			slog.Info("task queue is empty. waiting for next enqueue")
			cond.L.Lock()
			cond.Wait()
			cond.L.Unlock()
			return
		} else {
			slog.Error("task dispatch error", "error", err)
			return
		}
	}

	slog.Info("task retrieved", "type", task.Type, "description", task.Description)

	result, err := db.Exec("update tasks set status=? where id=? and status=?", statusRunning, task.ID, statusQueued)
	if err != nil {
		slog.Error("task update error", "error", err)
		return
	}

	n, err := result.RowsAffected()
	if err != nil {
		slog.Error("failed to get rows affected", "error", err)
		return
	}

	if n == 0 {
		slog.Info("failed to acquire task. retrying.")
		return
	}

	handler, ok := handlers[task.Type]
	if !ok {
		err = fallbackHandler(&task)
	} else {
		err = handler(&task)
	}

	status := statusFinished
	if err != nil {
		slog.Error("task handling error", "type", task.Type, "description", task.Description, "error", err)
		status = statusError
	}

	_, err = db.Exec("update tasks set status=? where id=?", status, task.ID)

	if err != nil {
		slog.Error("task update error", "error", err)
	}
}

func Handler(tasktype string, handler TaskHandler) error {
	if _, ok := handlers[tasktype]; ok {
		return errors.New("cannot reuse existing type: " + tasktype)
	}

	handlers[tasktype] = handler

	return nil
}

func Start() {
	for {
		Dispatch()
	}
}
