package dto

import (
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/task"
)

type TaskHistory struct {
	ID        int              `json:"id"`
	Name      string           `json:"name"`
	Result    model.TaskResult `json:"result"`
	RunAt     time.Time        `json:"run_at"`
	Error     string           `json:"error,omitempty"`
	Recurring bool             `json:"recurring"`
	Duration  time.Duration    `json:"duration"`
}

func TaskHistoryDTO(task *model.Task) TaskHistory {
	taskError := ""
	if task.Error != nil {
		taskError = task.Error.Error()
	}

	return TaskHistory{
		ID:        task.ID,
		Name:      task.Name,
		Result:    task.Result,
		RunAt:     task.RunAt,
		Error:     taskError,
		Recurring: task.Recurring,
		Duration:  task.Duration,
	}
}

type TaskHistoryFilter struct {
	Name   string
	Result *model.TaskResult
	Limit  int
	Offset int
}

type TaskStatus string

const (
	Running TaskStatus = "running"
	Waiting TaskStatus = "waiting"
)

type Task struct {
	ID         int              `json:"id"`
	Name       string           `json:"name"`
	Status     TaskStatus       `json:"status"`
	NextRun    time.Time        `json:"next_run"`
	Recurring  bool             `json:"recurring"`
	LastStatus model.TaskResult `json:"last_status,omitempty"`
	LastRun    *time.Time       `json:"last_run,omitempty"`
	LastError  string           `json:"last_error,omitempty"`
	Interval   *time.Duration   `json:"interval,omitempty"`
}

func TaskDTO(task task.Stat) Task {
	t := Task{
		ID:        task.ID,
		Name:      task.Name,
		Status:    TaskStatus(task.Status),
		NextRun:   task.NextRun,
		Recurring: task.Recurring,
	}

	if task.Recurring {
		t.LastStatus = model.TaskResult(task.LastStatus)
		t.LastRun = &task.LastRun
		if task.LastError != nil {
			t.LastError = task.LastError.Error()
		}
		t.Interval = &task.Interval
	}

	return t
}
