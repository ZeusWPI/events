// Package task provides an interface to schedule background delayed tasks
package task

import (
	"context"
	"time"

	"go.uber.org/zap"
)

// Now can be used if a one time task needs to be run immediately.
// It shouldn't be used with a recurring task
const Now = time.Duration(0)

// It's recommended to create a task with the NewTask function for added features
type Task interface {
	// This has to be unique for recurring tasks
	Name() string
	// Interval returns the time between executions.
	// For one time tasks it represents the amount of time to wait before executing
	// For immediate execution task.Now can be used
	Interval() time.Duration
	Func() func(context.Context) error
	Ctx() context.Context
}

type Status string

const (
	Waiting Status = "waiting"
	Running Status = "running"
)

type LastStatus string

const (
	Success LastStatus = "success"
	Failed  LastStatus = "failed"
)

// For one time tasks some fields are not used
// TODO: Check to split in Stat and StatRecurring
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
