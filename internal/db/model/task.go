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
	Resolved TaskResult = "resolved" // Task failed but has been marked as resolved by the user
)

type TaskType string

const (
	TaskRecurring TaskType = "recurring"
	TaskOnce      TaskType = "once"
)

type Task struct {
	// Task result
	ID       int // ID of the task result
	RunAt    time.Time
	Result   TaskResult
	Error    error
	Duration time.Duration

	// Task fields
	UID    string // Identifier of the task
	Name   string
	Active bool
	Type   TaskType
}

func TaskModel(task sqlc.Task, taskRun sqlc.TaskRun) *Task {
	var err error
	if taskRun.Error.Valid {
		err = errors.New(taskRun.Error.String)
	}

	return &Task{
		ID:       int(taskRun.ID),
		RunAt:    taskRun.RunAt.Time,
		Result:   TaskResult(taskRun.Result),
		Error:    err,
		Duration: time.Duration(taskRun.Duration),
		UID:      task.Uid,
		Name:     task.Name,
		Active:   task.Active,
		Type:     TaskType(task.Type),
	}
}

type TaskFilter struct {
	TaskUID string
	Result  *TaskResult
	Limit   int
	Offset  int
}
