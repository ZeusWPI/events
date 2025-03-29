// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: task.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const taskCreate = `-- name: TaskCreate :one
INSERT INTO task (name, result, run_at, error, recurring)
VALUES ($1, $2, $3, $4, $5)
RETURNING id
`

type TaskCreateParams struct {
	Name      string
	Result    pgtype.Text
	RunAt     pgtype.Timestamptz
	Error     pgtype.Text
	Recurring bool
}

func (q *Queries) TaskCreate(ctx context.Context, arg TaskCreateParams) (int32, error) {
	row := q.db.QueryRow(ctx, taskCreate,
		arg.Name,
		arg.Result,
		arg.RunAt,
		arg.Error,
		arg.Recurring,
	)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const taskGet = `-- name: TaskGet :many
SELECT id, name, result, run_at, error, recurring FROM task 
WHERE name ILIKE $1
ORDER BY run_at DESC
`

func (q *Queries) TaskGet(ctx context.Context, name string) ([]Task, error) {
	rows, err := q.db.Query(ctx, taskGet, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Result,
			&i.RunAt,
			&i.Error,
			&i.Recurring,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const taskGetAll = `-- name: TaskGetAll :many
SELECT id, name, result, run_at, error, recurring FROM task
ORDER BY run_at DESC
`

func (q *Queries) TaskGetAll(ctx context.Context) ([]Task, error) {
	rows, err := q.db.Query(ctx, taskGetAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Result,
			&i.RunAt,
			&i.Error,
			&i.Recurring,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
