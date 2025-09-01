package model

import (
	"errors"
	"time"

	"github.com/ZeusWPI/events/internal/db/sqlc"
)

type TaskResult string

const (
	Success  TaskResult = "success"
	Failed   TaskResult = "failed"
	Resolved TaskResult = "resolved"
)

type Task struct {
	ID        int
	Name      string
	Result    TaskResult
	RunAt     time.Time
	Error     error
	Recurring bool
	Duration  time.Duration
}

type TaskFilter struct {
	Name        string
	OnlyErrored bool
	Recurring   *bool
	Page        int
	Limit       int
}

func TaskModel(task sqlc.Task) *Task {
	var errTask error
	if task.Error.Valid {
		errTask = errors.New(task.Error.String)
	}

	return &Task{
		ID:        int(task.ID),
		Name:      task.Name,
		Result:    TaskResult(task.Result),
		RunAt:     task.RunAt.Time,
		Error:     errTask,
		Recurring: task.Recurring,
		Duration:  time.Duration(task.Duration.Microseconds * int64(time.Microsecond)),
	}
}
