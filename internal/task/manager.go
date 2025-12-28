package task

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/pkg/config"
	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

// Manager is the global single task manager instance
var Manager *manager

type job struct {
	task model.Task

	status Status
}

type jobRecurring struct {
	job

	interval   time.Duration
	lastStatus model.TaskResult
	lastError  error
}

type jobOnce struct {
	job
}

// Manager can be used to schedule one time or recurring tasks in the background
// It keeps logs inside the database.
// However it does not automatically reshedule tasks after an application reboot
type manager struct {
	scheduler gocron.Scheduler
	repo      repository.Task

	mu            sync.Mutex
	jobsRecurring map[string]jobRecurring
	jobsOnce      map[string]jobOnce

	development bool
	channelID   string
}

func newManager(repo repository.Repository) (*manager, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("failed to create new scheduler %w", err)
	}

	scheduler.Start()

	manager := &manager{
		scheduler:     scheduler,
		repo:          *repo.NewTask(),
		jobsRecurring: make(map[string]jobRecurring),
		jobsOnce:      make(map[string]jobOnce),
		development:   config.IsDev(),
		channelID:     config.GetString("task.channel"),
	}

	if err := manager.repo.SetInactiveRecurring(context.Background()); err != nil {
		return nil, err
	}

	return manager, nil
}

// AddRecurring adds a new recurring task to the manager.
// It immediately runs the task and then schedules it according to the interval.
// An unique uid is required.
// History logs (in the DB) for recurrent tasks are accessed by uid.
// If you change a recurring task's uid then all it's history will be lost (but still in the DB)
func (m *manager) AddRecurring(ctx context.Context, newTask Task) error {
	zap.S().Infof("Adding recurring task: %s | interval: %s", newTask.Name(), newTask.Interval())

	if _, ok := m.jobsRecurring[newTask.UID()]; ok {
		return fmt.Errorf("task %s already exists (uid: %s)", newTask.Name(), newTask.UID())
	}

	task, err := m.repo.GetByUID(ctx, newTask.UID())
	if err != nil {
		return err
	}
	if task != nil {
		// Pre-existing task
		// Update it
		task.Name = newTask.Name()
		task.Active = true
		if err := m.repo.Update(ctx, *task); err != nil {
			return err
		}
	} else {
		// New task
		// Let's create it
		task = &model.Task{
			UID:    newTask.UID(),
			Name:   newTask.Name(),
			Active: true,
			Type:   model.TaskRecurring,
		}
		if err := m.repo.Create(ctx, *task); err != nil {
			return err
		}
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Will immediately try to execute but it'll have to wait until the lock is released
	if _, err := m.scheduler.NewJob(
		gocron.DurationJob(newTask.Interval()),
		gocron.NewTask(m.wrapRecurring(newTask)),
		gocron.WithName(task.UID),
		gocron.WithContext(newTask.Ctx()),
		gocron.WithTags(task.UID),
		gocron.WithStartAt(gocron.WithStartImmediately()),
	); err != nil {
		return fmt.Errorf("failed to add task %+v | %w", *task, err)
	}

	m.jobsRecurring[task.UID] = jobRecurring{
		job: job{
			task:   *task,
			status: Waiting,
		},
		interval:   newTask.Interval(),
		lastStatus: model.Success,
		lastError:  nil,
	}

	return nil
}

// AddOnce adds a new one time task to the manager.
// It runs the tasks after the given interval and deletes it afterwards.
// An unique uid is required.
func (m *manager) AddOnce(newTask Task) error {
	zap.S().Infof("Adding one time task %s", newTask.Name())

	startTime := time.Now()
	startAtOption := gocron.OneTimeJobStartImmediately()
	if newTask.Interval() != Now {
		startTime = startTime.Add(newTask.Interval())
		startAtOption = gocron.OneTimeJobStartDateTime(startTime)
	}

	// No need to create an entry in the task table yet
	// It would only complicate it when it gets deleted
	// or cancelled because of a reboot
	task := &model.Task{
		UID:    newTask.UID(),
		Name:   newTask.Name(),
		Active: true,
		Type:   model.TaskOnce,
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// If startAtOption is set to immediately then it will immediately try to execute but it'll have to wait until the lock is released
	if _, err := m.scheduler.NewJob(
		gocron.OneTimeJob(startAtOption),
		gocron.NewTask(m.wrapOnce(newTask)),
		gocron.WithName(task.UID),
		gocron.WithContext(newTask.Ctx()),
		gocron.WithTags(task.UID),
	); err != nil {
		return fmt.Errorf("failed to add one time task %+v | %w", *task, err)
	}

	m.jobsOnce[task.UID] = jobOnce{
		job: job{
			task:   *task,
			status: Waiting,
		},
	}

	return nil
}

// RemoveByUID removes a scheduled task by its task UID
// Use RemoveByID to remove a task by it's job ID
func (m *manager) RemoveByUID(taskUID string) error {
	zap.S().Infof("Removing task by task uid %s", taskUID)

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.jobsOnce[taskUID]; ok {
		m.scheduler.RemoveByTags(taskUID)
		delete(m.jobsOnce, taskUID)

		return nil
	}

	if _, ok := m.jobsRecurring[taskUID]; ok {
		m.scheduler.RemoveByTags(taskUID)
		delete(m.jobsRecurring, taskUID)

		return nil
	}

	return fmt.Errorf("task with task uid %s not found", taskUID)
}

// RunByUID runs a pre existing task given a task UID.
func (m *manager) RunByUID(taskUID string) error {
	var job gocron.Job
	for _, j := range m.scheduler.Jobs() {
		if taskUID == j.Tags()[0] {
			job = j
			break
		}
	}
	if job == nil {
		return fmt.Errorf("task with uid %s not found", taskUID)
	}

	if err := job.RunNow(); err != nil {
		return fmt.Errorf("failed to run task with uid %s | %w", taskUID, err)
	}

	return nil
}

// Tasks returns all scheduled tasks
func (m *manager) Tasks() ([]Stat, error) {
	m.mu.Lock()
	jobs := m.scheduler.Jobs()
	jobsOnce := m.jobsOnce
	jobsRecurring := m.jobsRecurring
	m.mu.Unlock()

	stats := make([]Stat, 0, len(jobs))

	for _, job := range jobs {
		taskUID := job.Tags()[0]

		nextRun, err := job.NextRun()
		if err != nil {
			return nil, fmt.Errorf("get next run for task %s | %w", job.Name(), err)
		}

		if j, ok := jobsOnce[taskUID]; ok {
			stats = append(stats, Stat{
				TaskUID: j.task.UID,
				Name:    j.task.Name,
				Status:  j.status,
				NextRun: nextRun,
				Type:    j.task.Type,
			})
		} else if j, ok := jobsRecurring[taskUID]; ok {
			lastRun, err := job.LastRun()
			if err != nil {
				return nil, fmt.Errorf("get last run for task %s | %w", job.Name(), err)
			}

			stats = append(stats, Stat{
				TaskUID:    j.task.UID,
				Name:       j.task.Name,
				Status:     j.status,
				NextRun:    nextRun,
				Type:       j.task.Type,
				LastStatus: j.lastStatus,
				LastRun:    lastRun,
				LastError:  j.lastError,
				Interval:   j.interval,
			})
		}
	}

	slices.SortFunc(stats, func(a, b Stat) int { return int(a.NextRun.Sub(b.NextRun).Nanoseconds()) })

	return stats, nil
}

func (m *manager) wrapRecurring(task Task) func(context.Context) {
	return func(ctx context.Context) {
		m.mu.Lock()
		info, ok := m.jobsRecurring[task.UID()]
		if !ok {
			// Should not be possible
			m.mu.Unlock()
			zap.S().Errorf("Task %s not found during execution", task.Name())
			return
		}

		info.status = Running
		m.jobsRecurring[task.UID()] = info
		m.mu.Unlock()

		// Run task
		start := time.Now()
		err := task.Func()(ctx)
		end := time.Now()

		// Save result
		result := model.Success
		if err != nil {
			result = model.Failed
		}

		taskDB := &model.Task{
			UID:      task.UID(),
			RunAt:    time.Now(),
			Result:   result,
			Error:    err,
			Duration: end.Sub(start),
		}

		if errDB := m.repo.CreateRun(ctx, taskDB); errDB != nil {
			zap.S().Errorf("Failed to save recurring task result in database %+v | %v", *taskDB, err)
		}

		m.mu.Lock()
		defer m.mu.Unlock()

		info = m.jobsRecurring[task.UID()]
		info.status = Waiting
		info.lastStatus = result
		info.lastError = err
		m.jobsRecurring[task.UID()] = info
	}
}

func (m *manager) wrapOnce(task Task) func(context.Context) {
	return func(ctx context.Context) {
		m.mu.Lock()
		info, ok := m.jobsOnce[task.UID()]
		if !ok {
			// Should not be possible
			m.mu.Unlock()
			zap.S().Errorf("Task %s not found during execution", task.Name())
			return
		}

		info.status = Running
		m.jobsOnce[task.UID()] = info
		m.mu.Unlock()

		// Run task
		start := time.Now()
		taskErr := task.Func()(ctx)
		end := time.Now()

		// Save result
		taskDB, err := m.repo.GetByUID(ctx, task.UID())
		if err != nil {
			zap.S().Errorf("Failed to get pre-existing task %+v | %v", task, err)
			return
		}
		if taskDB == nil {
			// Will be the case 99% of the time
			// The only times it isn't true is if the one time task failed and an user reran it manually
			taskDB = &model.Task{
				UID:    task.UID(),
				Name:   task.Name(),
				Active: true,
				Type:   model.TaskOnce,
			}

			if err := m.repo.Create(ctx, *taskDB); err != nil {
				zap.S().Errorf("Failed to save one time task in database %+v | %v", *taskDB, err)
			}
		}

		result := model.Success
		if taskErr != nil {
			result = model.Failed
		}

		taskDB.RunAt = time.Now()
		taskDB.Result = result
		taskDB.Error = taskErr
		taskDB.Duration = end.Sub(start)

		if err := m.repo.CreateRun(ctx, taskDB); err != nil {
			zap.S().Errorf("Failed to save one time task run in database %+v | %v", *taskDB, err)
		}

		m.mu.Lock()
		defer m.mu.Unlock()

		delete(m.jobsOnce, task.UID())
	}
}
