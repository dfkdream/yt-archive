package taskq

import (
	"database/sql"
	"log"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

var defaultQueue *Queue

func SetDefaultQueue(q *Queue) {
	defaultQueue = q
}

func DefaultQueue() *Queue {
	if defaultQueue != nil {
		return defaultQueue
	}

	db, err := sql.Open("sqlite3", "file:taskq.db?_journal_mode=WAL")
	db.SetMaxOpenConns(1)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("opened sqlite database", "filename", "taskq.db")

	defaultQueue, err = New(db)
	if err != nil {
		log.Fatal(err)
	}

	return defaultQueue
}

func Enqueue(task *Task) error {
	return DefaultQueue().Enqueue(task)
}

func Handler(tasktype string, handler TaskHandler) error {
	return DefaultQueue().Handler(tasktype, handler)
}

func Start() {
	DefaultQueue().Start()
}

func Close() error {
	if defaultQueue == nil {
		slog.Info("default task queue is nil")
		return nil
	}

	slog.Info("closing database for default task queue")
	return defaultQueue.db.Close()
}
