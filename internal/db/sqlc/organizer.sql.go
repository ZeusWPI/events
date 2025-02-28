// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: organizer.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const organizerCreate = `-- name: OrganizerCreate :one
INSERT INTO organizer (event, board)
VALUES ($1, $2)
RETURNING id
`

type OrganizerCreateParams struct {
	Event int32
	Board int32
}

func (q *Queries) OrganizerCreate(ctx context.Context, arg OrganizerCreateParams) (int32, error) {
	row := q.db.QueryRow(ctx, organizerCreate, arg.Event, arg.Board)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const organizerDelete = `-- name: OrganizerDelete :exec
DELETE FROM organizer 
WHERE id = $1
`

func (q *Queries) OrganizerDelete(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, organizerDelete, id)
	return err
}

const organizerGetByYearWithBoard = `-- name: OrganizerGetByYearWithBoard :many
SELECT o.id, event, board, b.id, member, year, role, created_at, updated_at, m.id, name, username FROM organizer o 
INNER JOIN board b ON b.id = o.board 
INNER JOIN member m ON m.id = b.member 
WHERE b.year = $1
`

type OrganizerGetByYearWithBoardRow struct {
	ID        int32
	Event     int32
	Board     int32
	ID_2      int32
	Member    int32
	Year      int32
	Role      string
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	ID_3      int32
	Name      string
	Username  pgtype.Text
}

func (q *Queries) OrganizerGetByYearWithBoard(ctx context.Context, year int32) ([]OrganizerGetByYearWithBoardRow, error) {
	rows, err := q.db.Query(ctx, organizerGetByYearWithBoard, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []OrganizerGetByYearWithBoardRow
	for rows.Next() {
		var i OrganizerGetByYearWithBoardRow
		if err := rows.Scan(
			&i.ID,
			&i.Event,
			&i.Board,
			&i.ID_2,
			&i.Member,
			&i.Year,
			&i.Role,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ID_3,
			&i.Name,
			&i.Username,
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
