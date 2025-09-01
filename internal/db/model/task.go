package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/ZeusWPI/events/internal/db/sqlc"
)

type TaskResult string

const (
	Success  TaskResult = "success"
	Failed   TaskResult = "failed"
	Resolved TaskResult = "resolved"
)

func TaskResultModel(result string) (TaskResult, error) {
	switch result {
	case string(Success), string(Failed), string(Resolved):
		return TaskResult(result), nil
	default:
		return "", fmt.Errorf("invalid task result %s", result)
	}
}

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
	Name   string
	Result *TaskResult
	Limit  int
	Offset int
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
