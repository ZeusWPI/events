// Package task provides a manager to schedule one time and recurring tasks
package task

import (
	"context"
	"time"

	"go.uber.org/zap"
)

// Now can be used if a one time task needs to be run immediately.
// It shouldn't be used with a recurring task
const Now = time.Duration(0)

// Task represents a task.
// It's recommended to create a task with the NewTask function for added features
type Task interface {
	// Name returns the name of the Task.
	// This has to be unique for recurring tasks
	Name() string
	// Interval returns the time between executions.
	// For one time tasks it represents the amount of time to wait before executing
	// For immediate execution task.Now can be used
	Interval() time.Duration
	Func() func(context.Context) error
	Ctx() context.Context
}

// Status represents the possible statuses of a task
type Status string

const (
	// Waiting indicates a task is waiting to be executed
	Waiting Status = "waiting"
	// Running indicates a task is being executed
	Running Status = "running"
)

// LastStatus represents the possible statuses of the previous execution of a task
type LastStatus string

const (
	// Success indicates a task exited without any errors
	Success LastStatus = "success"
	// Failed indicates a task exited with an error
	Failed LastStatus = "failed"
)

// Stat represents a task with some statistics .
// For one time tasks some fields are not used
type Stat struct {
	ID         int
	Name       string
	Status     Status
	NextRun    time.Time
	Recurring  bool
	LastStatus LastStatus    // Not used with one time tasks
	LastRun    time.Time     // Not used with one time tasks
	LastError  error         // Not used with one time tasks
	Interval   time.Duration // Not used with one time tasks
}

type internalTask struct {
	name     string
	interval time.Duration
	fn       func(context.Context) error
	ctx      context.Context
}

var _ Task = (*internalTask)(nil)

// NewTask creates a new task
// It supports an optional context, if none is given the background context is used
// Logs (info level) when a task starts
// Logs (error level) any error that occurs during the task execution
func NewTask(name string, interval time.Duration, fn func(context.Context) error, ctx ...context.Context) Task {
	c := context.Background()
	if len(ctx) > 0 {
		c = ctx[0]
	}

	return &internalTask{
		name:     name,
		interval: interval,
		fn:       fn,
		ctx:      c,
	}
}

func (t *internalTask) Name() string {
	return t.name
}

func (t *internalTask) Interval() time.Duration {
	return t.interval
}

func (t *internalTask) Func() func(context.Context) error {
	return func(ctx context.Context) error {
		zap.S().Infof("Running task %s", t.name)

		if err := t.fn(ctx); err != nil {
			zap.S().Errorf("Task %s failed | %v", t.name, err)
			return err
		}

		return nil
	}
}

func (t *internalTask) Ctx() context.Context {
	return t.ctx
}
