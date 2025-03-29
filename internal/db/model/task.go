package model

import "time"

// TaskResult is the different results a task can have
type TaskResult string

const (
	// Success indicates a task succeeded
	Success TaskResult = "success"
	// Failed indicates a task failed
	Failed TaskResult = "failed"
)

// Task represents a task
type Task struct {
	ID        int
	Name      string
	Result    TaskResult
	RunAt     time.Time
	Error     error
	Recurring bool
}

// TaskFilter allows for filtering tasks
type TaskFilter struct {
	Name        string
	OnlyErrored bool
	Recurring   *bool
	Page        int
	Limit       int
}
