package dto

import (
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/task"
)

type TaskHistoryStatus string

const (
	Success TaskHistoryStatus = "Success"
	Failed  TaskHistoryStatus = "Failed"
)

type TaskHistory struct {
	ID        int               `json:"id"`
	Name      string            `json:"name"`
	Result    TaskHistoryStatus `json:"result"`
	RunAt     time.Time         `json:"run_at"`
	Error     string            `json:"error,omitempty"`
	Recurring bool              `json:"recurring"`
}

func TaskHistoryDTO(task *model.Task) TaskHistory {
	taskError := ""
	if task.Error != nil {
		taskError = task.Error.Error()
	}

	return TaskHistory{
		ID:        task.ID,
		Name:      task.Name,
		Result:    TaskHistoryStatus(task.Result),
		RunAt:     task.RunAt,
		Error:     taskError,
		Recurring: task.Recurring,
	}
}

type TaskHistoryFilter struct {
	Name        string
	OnlyErrored bool
	Recurring   *bool
	Page        int
	Limit       int
}

type TaskStatus string

const (
	Running TaskStatus = "Running"
	Waiting TaskStatus = "Waiting"
)

type Task struct {
	ID         int               `json:"id"`
	Name       string            `json:"name"`
	Status     TaskStatus        `json:"status"`
	NextRun    time.Time         `json:"next_run"`
	Recurring  bool              `json:"recurring"`
	LastStatus TaskHistoryStatus `json:"last_status,omitempty"`
	LastRun    *time.Time        `json:"last_run,omitempty"`
	LastError  string            `json:"last_error,omitempty"`
	Interval   *time.Duration    `json:"interval,omitempty"`
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
		t.LastStatus = TaskHistoryStatus(task.LastStatus)
		t.LastRun = &task.LastRun
		t.LastError = ""
		if task.LastError != nil {
			t.LastError = task.LastError.Error()
		}
		t.Interval = &task.Interval
	}

	return t
}
