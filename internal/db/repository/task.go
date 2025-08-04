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

func (t *Task) GetFiltered(ctx context.Context, filters model.TaskFilter) ([]*model.Task, error) {
	tasks, err := t.repo.queries(ctx).TaskGetAll(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get all tasks %w", err)
	}

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
	end := min(start+filters.Limit, len(filtered))

	return utils.SliceMap(filtered[start:end], model.TaskModel), nil
}

func (t *Task) Create(ctx context.Context, task *model.Task) error {
	errTask := pgtype.Text{Valid: false}
	if task.Error != nil {
		errTask.String = task.Error.Error()
		errTask.Valid = true
	}

	id, err := t.repo.queries(ctx).TaskCreate(ctx, sqlc.TaskCreateParams{
		Name:      task.Name,
		Result:    pgtype.Text{String: string(task.Result), Valid: true},
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
