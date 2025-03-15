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

const organizerDeleteByBoardEvent = `-- name: OrganizerDeleteByBoardEvent :exec
DELETE FROM organizer 
WHERE board = $1 AND event = $2
`

type OrganizerDeleteByBoardEventParams struct {
	Board int32
	Event int32
}

func (q *Queries) OrganizerDeleteByBoardEvent(ctx context.Context, arg OrganizerDeleteByBoardEventParams) error {
	_, err := q.db.Exec(ctx, organizerDeleteByBoardEvent, arg.Board, arg.Event)
	return err
}

const organizerGetByYearWithBoard = `-- name: OrganizerGetByYearWithBoard :many
SELECT o.id, event, board, b.id, member, year, role, created_at, updated_at, m.id, name, username, zauth_id FROM organizer o 
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
	ZauthID   pgtype.Int4
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
			&i.ZauthID,
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
