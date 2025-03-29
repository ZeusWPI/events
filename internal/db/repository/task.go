package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/sqlc"
	"github.com/ZeusWPI/events/pkg/util"
	"github.com/jackc/pgx/v5/pgtype"
)

// Task provides all model.Task related database operations
type Task interface {
	GetByName(context.Context, string) ([]*model.Task, error)
	GetFiltered(context.Context, model.TaskFilter) ([]*model.Task, error)
	Save(context.Context, *model.Task) error
}

type taskRepo struct {
	repo Repository
}

// Interface compliance
var _ Task = (*taskRepo)(nil)

func (r *taskRepo) GetByName(ctx context.Context, name string) ([]*model.Task, error) {
	tasks, err := r.repo.queries(ctx).TaskGet(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get task %s | %w", name, err)
	}

	return util.SliceMap(tasks, r.convert), nil
}

func (r *taskRepo) GetFiltered(ctx context.Context, filters model.TaskFilter) ([]*model.Task, error) {
	tasks, err := r.repo.queries(ctx).TaskGetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all tasks %w", err)
	}

	// TODO: Do filtering & pagination in SQL

	filtered := tasks[:0]

	// Filtering
	for _, task := range tasks {
		if filters.Name != "" && task.Name != filters.Name {
			continue
		}
		if filters.OnlyErrored && !task.Error.Valid {
			continue
		}
		if filters.Recurring != nil && task.Recurring != *filters.Recurring {
			continue
		}

		filtered = append(filtered, task)
	}

	// Pagination
	start := (filters.Page - 1) * filters.Limit
	if start >= len(filtered) {
		return []*model.Task{}, nil
	}
	if start < 0 {
		start = 0
	}
	end := start + filters.Limit
	if end > len(filtered) {
		end = len(filtered)
	}

	return util.SliceMap(filtered[start:end], r.convert), nil
}

func (r *taskRepo) Save(ctx context.Context, t *model.Task) error {
	errTask := pgtype.Text{Valid: false}
	if t.Error != nil {
		errTask.String = t.Error.Error()
		errTask.Valid = true
	}

	id, err := r.repo.queries(ctx).TaskCreate(ctx, sqlc.TaskCreateParams{
		Name:      t.Name,
		Result:    pgtype.Text{String: string(t.Result), Valid: true},
		RunAt:     pgtype.Timestamptz{Time: t.RunAt, Valid: true},
		Error:     errTask,
		Recurring: t.Recurring,
	})
	if err != nil {
		return fmt.Errorf("unable to save task %+v | %w", *t, err)
	}

	t.ID = int(id)

	return nil
}

func (r *taskRepo) convert(task sqlc.Task) *model.Task {
	var errTask error
	if task.Error.Valid {
		errTask = errors.New(task.Error.String)
	}

	return &model.Task{
		ID:        int(task.ID),
		Name:      task.Name,
		Result:    model.TaskResult(task.Result.String),
		RunAt:     task.RunAt.Time,
		Error:     errTask,
		Recurring: task.Recurring,
	}
}
