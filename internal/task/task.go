// Package task provides an interface to schedule background delayed tasks
package task

import (
	"context"
	"fmt"
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/pkg/mattermost"
	"go.uber.org/zap"
)

// Init intializes the global task manager instance
func Init(repo repository.Repository) error {
	manager, err := newManager(repo)
	if err != nil {
		return err
	}

	Manager = manager

	return nil
}

// Now can be used if a one time task needs to be run immediately.
// It shouldn't be used with a recurring task
const Now = time.Duration(0)

// Task is the interface to whcih a task should adhere to
// You can manually implement all methods are make use of the `NewTask` function
// which will automatically add some logging
type Task interface {
	// UID is an unique identifier for a check
	// History is kept by linking the UID's of tasks
	// Changing the UID will make you lose all the task history
	// Changing the frontend name can be done with the Name() function
	UID() string
	// Name is an user friendly task name
	// You can change this as much as you like
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

// Stat contains the information about a current running or scheduled task
// Some fields are populated depending on the type (recurring or one-time)
type Stat struct {
	TaskUID    string
	Name       string
	Status     Status
	NextRun    time.Time
	Type       model.TaskType
	LastStatus model.TaskResult // Not used with one time tasks
	LastRun    time.Time        // Not used with one time tasks
	LastError  error            // Not used with one time tasks
	Interval   time.Duration    // Not used with one time tasks
}

type internalTask struct {
	uid      string
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
func NewTask(uid string, name string, interval time.Duration, fn func(context.Context) error, ctx ...context.Context) Task {
	c := context.Background()
	if len(ctx) > 0 {
		c = ctx[0]
	}

	return &internalTask{
		uid:      uid,
		name:     name,
		interval: interval,
		fn:       fn,
		ctx:      c,
	}
}

func (t *internalTask) UID() string {
	return t.uid
}

func (t *internalTask) Name() string {
	return t.name
}

func (t *internalTask) Interval() time.Duration {
	return t.interval
}

func (t *internalTask) Func() func(context.Context) error {
	return func(ctx context.Context) error {
		zap.S().Infof("Task running %s", t.name)

		if err := t.fn(ctx); err != nil {
			zap.S().Errorf("Task %s failed | %v", t.name, err)

			message := fmt.Sprintf("Task %s failed\n%v", t.name, err)
			if Manager.development {
				zap.S().Infof("Mock task failed mattermost message: \n%s", message)
			} else {
				if err := mattermost.C.SendMessage(ctx, mattermost.Message{
					ChannelID: Manager.channelID,
					Message:   message,
				}); err != nil {
					zap.S().Errorf("Send mattermost message failed %v", err)
				}
			}

			return err
		}

		zap.S().Infof("Task finished %s", t.name)
		return nil
	}
}

func (t *internalTask) Ctx() context.Context {
	return t.ctx
}
