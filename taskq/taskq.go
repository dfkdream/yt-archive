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

type Task struct {
	ID          uuid.UUID
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
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("opened sqlite database", "filename", "taskq.db")

	_, err = db.Exec("create table if not exists tasks_queued (id blob primary key, priority integer, type text, description text, payload text)")
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("created table tasks_queued")

	_, err = db.Exec("create table if not exists tasks_finished (id blob primary key, type text, description text, payload text)")
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("created table tasks_finished")

	handlers = make(map[string]TaskHandler)
}

func NewTask(priority int, tasktype, description string, payload string) (*Task, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return &Task{
		ID:          uuid,
		Priority:    priority,
		Type:        tasktype,
		Description: description,
		Payload:     payload,
	}, nil
}

func Enqueue(task *Task) error {
	_, err := db.Exec("insert into tasks_queued (id, priority, type, description, payload) values (?, ?, ?, ?, ?)",
		task.ID, task.Priority, task.Type, task.Description, task.Payload)

	if err == nil {
		slog.Info("task enqueued", "type", task.Type, "description", task.Description)
		cond.L.Lock()
		cond.Signal()
		cond.L.Unlock()
	}

	return err
}

func Dispatch() {
	row := db.QueryRow("select id, priority, type, description, payload from tasks_queued order by priority desc, id asc limit 1")

	var task Task
	err := row.Scan(&task.ID, &task.Priority, &task.Type, &task.Description, &task.Payload)
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

	handler, ok := handlers[task.Type]
	if !ok {
		err = fallbackHandler(&task)
	} else {
		err = handler(&task)
	}

	if err != nil {
		slog.Error("task handling error", "type", task.Type, "description", task.Description, "error", err)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		slog.Error("transaction begin error", "error", err)
	}

	_, err = tx.Exec("insert into tasks_finished (id, type, description, payload) values (?, ?, ?, ?)",
		task.ID, task.Type, task.Description, task.Payload)

	if err != nil {
		tx.Rollback()
		slog.Error("insert error", "error", err)
		return
	}

	_, err = tx.Exec("delete from tasks_queued where id=?", task.ID)
	if err != nil {
		tx.Rollback()
		slog.Error("delete error", "error", err)
		return
	}

	err = tx.Commit()
	if err != nil {
		slog.Error("commit error", "error", err)
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
