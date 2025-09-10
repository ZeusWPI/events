package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/sqlc"
	"github.com/ZeusWPI/events/pkg/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

type Task struct {
	repo Repository
}

func (r *Repository) NewTask() *Task {
	return &Task{
		repo: *r,
	}
}

// GetByUID returns the task given an uid without any runs
func (t *Task) GetByUID(ctx context.Context, taskUID string) (*model.Task, error) {
	task, err := t.repo.queries(ctx).TaskGetByUID(ctx, taskUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get task by uid %s | %w", taskUID, err)
	}

	return model.TaskModel(task, sqlc.TaskRun{}), nil
}

func (t *Task) GetByRunID(ctx context.Context, runID int) (*model.Task, error) {
	task, err := t.repo.queries(ctx).TaskRunGet(ctx, int32(runID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get task run %d | %w", runID, err)
	}

	return model.TaskModel(task.Task, task.TaskRun), err
}

func (t *Task) GetFiltered(ctx context.Context, filter model.TaskFilter) ([]*model.Task, error) {
	// A default value is required for sql
	// If no result filter is set then it will ignored anyway
	result := model.Success
	if filter.Result != nil {
		result = *filter.Result
	}

	params := sqlc.TaskGetFilteredParams{
		Uid:           filter.TaskUID,
		FilterTaskUid: filter.TaskUID != "",
		Result:        sqlc.TaskResult(result),
		FilterResult:  filter.Result != nil,
		Limit:         int32(filter.Limit),
		Offset:        int32(filter.Offset),
	}

	tasks, err := t.repo.queries(ctx).TaskGetFiltered(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get filtered tasks %+v | %w", filter, err)
	}

	return utils.SliceMap(tasks, func(task sqlc.TaskGetFilteredRow) *model.Task {
		return model.TaskModel(task.Task, task.TaskRun)
	}), nil
}

func (t *Task) Create(ctx context.Context, task model.Task) error {
	if err := t.repo.queries(ctx).TaskCreate(ctx, sqlc.TaskCreateParams{
		Uid:    task.UID,
		Name:   task.Name,
		Active: task.Active,
		Type:   sqlc.TaskType(task.Type),
	}); err != nil {
		return fmt.Errorf("create task %+v | %w", t, err)
	}

	return nil
}

func (t *Task) CreateRun(ctx context.Context, task *model.Task) error {
	errStr := ""
	if task.Error != nil {
		errStr = task.Error.Error()
	}

	id, err := t.repo.queries(ctx).TaskRunCreate(ctx, sqlc.TaskRunCreateParams{
		TaskUid:  task.UID,
		RunAt:    pgtype.Timestamptz{Time: task.RunAt, Valid: !task.RunAt.IsZero()},
		Result:   sqlc.TaskResult(task.Result),
		Error:    pgtype.Text{String: errStr, Valid: errStr != ""},
		Duration: task.Duration.Nanoseconds(),
	})
	if err != nil {
		return fmt.Errorf("create task run %+v | %w", *task, err)
	}

	task.ID = int(id)

	return nil
}

func (t *Task) Update(ctx context.Context, task model.Task) error {
	if err := t.repo.queries(ctx).TaskUpdate(ctx, sqlc.TaskUpdateParams{
		Uid:    task.UID,
		Name:   task.Name,
		Active: task.Active,
	}); err != nil {
		return fmt.Errorf("update task %+v | %w", task, err)
	}

	return nil
}

func (t *Task) SetInactiveRecurring(ctx context.Context) error {
	if err := t.repo.queries(ctx).TaskSetInactiveRecurring(ctx); err != nil {
		return fmt.Errorf("set recurring tasks to inactive %w", err)
	}

	return nil
}

func (t *Task) RunResolve(ctx context.Context, runID int) error {
	if err := t.repo.queries(ctx).TaskRunResolve(ctx, int32(runID)); err != nil {
		return fmt.Errorf("task run resolve %d | %w", runID, err)
	}

	return nil
}
