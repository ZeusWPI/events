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

func (t *Task) Get(ctx context.Context, taskID int) (*model.Task, error) {
	task, err := t.repo.queries(ctx).TaskGet(ctx, int32(taskID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get task %d | %w", taskID, err)
	}

	return model.TaskModel(task), err
}

func (t *Task) GetFiltered(ctx context.Context, filter model.TaskFilter) ([]*model.Task, error) {
	result := sqlc.TaskResult(model.Success)
	if filter.Result != nil {
		result = sqlc.TaskResult(*filter.Result)
	}

	params := sqlc.TaskGetFilteredParams{
		Name:         filter.Name,
		FilterName:   filter.Name != "",
		Result:       result,
		FilterResult: filter.Result != nil,
		Limit:        int32(filter.Limit),
		Offset:       int32(filter.Offset),
	}

	tasks, err := t.repo.queries(ctx).TaskGetFiltered(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get filtered tasks %+v | %w", filter, err)
	}

	return utils.SliceMap(tasks, model.TaskModel), nil
}

func (t *Task) Create(ctx context.Context, task *model.Task) error {
	errTask := pgtype.Text{Valid: false}
	if task.Error != nil {
		errTask.String = task.Error.Error()
		errTask.Valid = true
	}

	id, err := t.repo.queries(ctx).TaskCreate(ctx, sqlc.TaskCreateParams{
		Name:      task.Name,
		Result:    sqlc.TaskResult(task.Result),
		RunAt:     pgtype.Timestamptz{Time: task.RunAt, Valid: true},
		Error:     errTask,
		Recurring: task.Recurring,
		Duration:  pgtype.Interval{Microseconds: task.Duration.Microseconds(), Valid: true},
	})
	if err != nil {
		return fmt.Errorf("create task %+v | %w", *t, err)
	}

	task.ID = int(id)

	return nil
}

func (t *Task) UpdateResult(ctx context.Context, task model.Task) error {
	if err := t.repo.queries(ctx).TaskUpdateResult(ctx, sqlc.TaskUpdateResultParams{
		ID:     int32(task.ID),
		Result: sqlc.TaskResult(task.Result),
	}); err != nil {
		return fmt.Errorf("task update result %+v | %w", task, err)
	}

	return nil
}
