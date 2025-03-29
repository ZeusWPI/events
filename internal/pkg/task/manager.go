package task

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/pkg/util"
	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

type job struct {
	id        int
	name      string
	status    Status
	recurring bool
}

type jobRecurring struct {
	job

	interval   time.Duration
	lastStatus LastStatus
	lastError  error
}

type jobOnce struct {
	job
}

// Manager manages all tasks.
// It keeps a logs inside the database.
// However it does not automatically reshedule tasks after a application reboot
type Manager struct {
	scheduler gocron.Scheduler
	repo      repository.Task

	mu            sync.Mutex
	jobID         int
	jobsRecurring map[int]jobRecurring
	jobsOnce      map[int]jobOnce
}

// NewManager creates a new Manager
func NewManager(repo repository.Repository) (*Manager, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("failed to create new scheduler %w", err)
	}

	scheduler.Start()

	return &Manager{
		scheduler:     scheduler,
		repo:          repo.NewTask(),
		jobID:         1,
		jobsRecurring: make(map[int]jobRecurring),
		jobsOnce:      make(map[int]jobOnce),
	}, nil
}

// Add adds a new recurring task to the manager.
// It immediately runs the task and then schedules it according to the interval.
// Recurring tasks are required to have an unique name.
// History logs (in the DB) for recurrent tasks are accessed by name.
// If you change a recurring task's name then all it's history will be lost (but still in the DB)
func (m *Manager) Add(task Task) error {
	zap.S().Debugf("Adding recurring task %s", task.Name())

	for _, v := range m.jobsRecurring {
		if v.name == task.Name() {
			return fmt.Errorf("task %s already exists", task.Name())
		}
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Will immediately try to execute but it'll have to wait until the lock is released
	if _, err := m.scheduler.NewJob(
		gocron.DurationJob(task.Interval()),
		gocron.NewTask(m.wrapRecurring(m.jobID, task)),
		gocron.WithName(task.Name()),
		gocron.WithContext(task.Ctx()),
		gocron.WithTags(strconv.Itoa(m.jobID)),
		gocron.WithStartAt(gocron.WithStartImmediately()),
	); err != nil {
		return fmt.Errorf("failed to add task %s | %w", task.Name(), err)
	}

	m.jobsRecurring[m.jobID] = jobRecurring{
		job: job{
			id:        m.jobID,
			name:      task.Name(),
			status:    Waiting,
			recurring: true,
		},
		interval:   task.Interval(),
		lastStatus: Success,
		lastError:  nil,
	}
	m.jobID++

	return nil
}

// AddOnce adds a new one time task to the manager.
// It runs the tasks after the given interval and deletes it afterwards.
func (m *Manager) AddOnce(task Task) error {
	zap.S().Debugf("Adding one time task %s", task.Name())

	startTime := time.Now()
	startAtOption := gocron.OneTimeJobStartImmediately()
	if task.Interval() != Now {
		startTime = startTime.Add(task.Interval())
		startAtOption = gocron.OneTimeJobStartDateTime(startTime)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// If startAtOption is set to immediately then it will immediately try to execute but it'll have to wait until the lock is released
	if _, err := m.scheduler.NewJob(
		gocron.OneTimeJob(startAtOption),
		gocron.NewTask(m.wrapOnce(m.jobID, task)),
		gocron.WithName(task.Name()),
		gocron.WithContext(task.Ctx()),
		gocron.WithTags(strconv.Itoa(m.jobID)),
	); err != nil {
		return fmt.Errorf("failed to add one time task %s | %w", task.Name(), err)
	}

	m.jobsOnce[m.jobID] = jobOnce{
		job: job{
			id:        m.jobID,
			name:      task.Name(),
			status:    Waiting,
			recurring: false,
		},
	}
	m.jobID++

	return nil
}

// Run runs a pre existing task given an id.
func (m *Manager) Run(id int) error {
	return m.run(strconv.Itoa(id))
}

// RunByName runs a pre existing task given a name.
// This is only relevant for recurring tasks
func (m *Manager) RunByName(name string) error {
	m.mu.Lock()

	var job *jobRecurring
	for _, v := range m.jobsRecurring {
		if v.name == name {
			job = &v
			break
		}
	}

	m.mu.Unlock()
	if job == nil {
		return fmt.Errorf("task %s not found", name)
	}

	return m.run(strconv.Itoa(job.id))
}

// Tasks returns all scheduled tasks
func (m *Manager) Tasks() ([]Stat, error) {
	m.mu.Lock()
	jobs := m.scheduler.Jobs()
	jobsOnce := m.jobsOnce
	jobsRecurring := m.jobsRecurring
	m.mu.Unlock()

	stats := make([]Stat, 0, len(jobs))
	var errs []error

	for _, job := range jobs {
		id, err := strconv.Atoi(job.Tags()[0])
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to convert job id %s tag to int %s | %w", job.Tags()[0], job.Name(), err))
			continue
		}

		nextRun, err := job.NextRun()
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to get next run for task %s | %w", job.Name(), err))
			continue
		}

		if j, found := jobsOnce[id]; found {
			stats = append(stats, Stat{
				ID:        j.id,
				Name:      j.name,
				Status:    j.status,
				NextRun:   nextRun,
				Recurring: false,
			})
		} else if j, found := jobsRecurring[id]; found {
			lastRun, err := job.LastRun()
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to get last run for task %s | %w", job.Name(), err))
				continue
			}

			stats = append(stats, Stat{
				ID:         j.id,
				Name:       j.name,
				Status:     j.status,
				NextRun:    nextRun,
				Recurring:  true,
				LastStatus: j.lastStatus,
				LastRun:    lastRun,
				LastError:  j.lastError,
				Interval:   j.interval,
			})
		}
	}

	if errs != nil {
		return nil, errors.Join(errs...)
	}

	slices.SortFunc(stats, func(a, b Stat) int { return int(a.NextRun.Sub(b.NextRun).Milliseconds()) })

	return stats, nil
}

func (m *Manager) wrapRecurring(id int, task Task) func(context.Context) {
	return func(ctx context.Context) {
		m.mu.Lock()
		info, ok := m.jobsRecurring[id]
		if !ok {
			// Should not be possible
			m.mu.Unlock()
			zap.S().Errorf("Task %s not found during execution", task.Name())
			return
		}
		info.status = Running
		m.jobsRecurring[id] = info
		m.mu.Unlock()

		// Run task
		err := task.Func()(ctx)

		// Save result
		task := &model.Task{
			Name:      task.Name(),
			Result:    util.If(err == nil, model.Success, model.Failed),
			RunAt:     time.Now(),
			Error:     err,
			Recurring: true,
		}
		if errDB := m.repo.Save(ctx, task); errDB != nil {
			zap.S().Errorf("failed to save recurring task result in database %+v | %v", *task, err)
		}

		m.mu.Lock()
		defer m.mu.Unlock()

		info = m.jobsRecurring[id]
		info.status = Waiting
		info.lastStatus = util.If(err == nil, Success, Failed)
		info.lastError = err
		m.jobsRecurring[id] = info
	}
}

func (m *Manager) wrapOnce(id int, task Task) func(context.Context) {
	return func(ctx context.Context) {
		m.mu.Lock()
		info, ok := m.jobsOnce[id]
		if !ok {
			// Should not be possible
			m.mu.Unlock()
			zap.S().Errorf("Task %s not found during execution", task.Name())
			return
		}
		info.status = Running
		m.jobsOnce[id] = info
		m.mu.Unlock()

		// Run task
		err := task.Func()(ctx)

		// Save result
		task := &model.Task{
			Name:      task.Name(),
			Result:    util.If(err == nil, model.Success, model.Failed),
			RunAt:     time.Now(),
			Error:     err,
			Recurring: false,
		}
		if err := m.repo.Save(ctx, task); err != nil {
			zap.S().Errorf("failed to save one time task result in database %+v | %v", *task, err)
		}

		m.mu.Lock()
		defer m.mu.Unlock()

		delete(m.jobsOnce, id)
	}
}

func (m *Manager) run(id string) error {
	m.mu.Lock()

	var job gocron.Job
	for _, j := range m.scheduler.Jobs() {
		if id == j.Tags()[0] {
			job = j
			break
		}
	}
	m.mu.Unlock()
	if job == nil {
		return fmt.Errorf("task with id %s not found", id)
	}

	if err := job.RunNow(); err != nil {
		return fmt.Errorf("failed to run task with id %s | %w", id, err)
	}

	return nil
}
