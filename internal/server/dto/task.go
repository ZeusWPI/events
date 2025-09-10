package dto

import (
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/task"
)

type TaskHistory struct {
	ID       int              `json:"id"`
	Name     string           `json:"name"`
	Result   model.TaskResult `json:"result"`
	RunAt    time.Time        `json:"run_at"`
	Error    string           `json:"error,omitempty"`
	Type     model.TaskType   `json:"type"`
	Duration time.Duration    `json:"duration"`
}

func TaskHistoryDTO(task *model.Task) TaskHistory {
	taskError := ""
	if task.Error != nil {
		taskError = task.Error.Error()
	}

	return TaskHistory{
		ID:       task.ID,
		Name:     task.Name,
		Result:   task.Result,
		RunAt:    task.RunAt,
		Error:    taskError,
		Type:     task.Type,
		Duration: task.Duration,
	}
}

type Task struct {
	TaskUID    string           `json:"uid"`
	Name       string           `json:"name"`
	Status     task.Status      `json:"status"`
	NextRun    time.Time        `json:"next_run"`
	Type       model.TaskType   `json:"type"`
	LastStatus model.TaskResult `json:"last_status,omitempty"`
	LastRun    *time.Time       `json:"last_run,omitzero"`
	LastError  string           `json:"last_error,omitempty"`
	Interval   *time.Duration   `json:"interval,omitzero"`
}

func TaskDTO(task task.Stat) Task {
	t := Task{
		TaskUID: task.TaskUID,
		Name:    task.Name,
		Status:  task.Status,
		NextRun: task.NextRun,
		Type:    task.Type,
	}

	if task.Type == model.TaskRecurring {
		t.LastStatus = task.LastStatus
		t.LastRun = &task.LastRun
		if task.LastError != nil {
			t.LastError = task.LastError.Error()
		}
		t.Interval = &task.Interval
	}

	return t
}

type TaskFilter struct {
	TaskUID string
	Result  *model.TaskResult
	Limit   int
	Offset  int
}
