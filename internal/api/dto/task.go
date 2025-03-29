package dto

import (
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/pkg/task"
)

// TaskHistoryStatus are the possible results a TaskHistory
type TaskHistoryStatus string

const (
	// Success indicates the task succeeded
	Success TaskHistoryStatus = "Success"
	// Failed indicates the task returned an error
	Failed TaskHistoryStatus = "Failed"
)

// TaskHistory is the data transferable object version of the model Task
type TaskHistory struct {
	ID        int               `json:"id"`
	Name      string            `json:"name"`
	Result    TaskHistoryStatus `json:"result"`
	RunAt     time.Time         `json:"run_at"`
	Error     string            `json:"error,omitempty"`
	Recurring bool              `json:"recurring"`
}

// TaskHistoryDTO converts a model Task to a DTO TaskHistory
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

// TaskHistoryFilter is the filter used for task histories
type TaskHistoryFilter struct {
	Name        string
	OnlyErrored bool
	Recurring   *bool
	Page        int
	Limit       int
}

// TaskStatus are the different statuses a Task can have
type TaskStatus string

const (
	// Running indicates the task is currently running
	Running TaskStatus = "Running"
	// Waiting indicates the task is currently waiting to be executed
	Waiting TaskStatus = "Waiting"
)

// Task is the data transferable object version of a task.Stat returned by the task manager
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

// TaskDTO converts a task.Stat to a DTO Task
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
