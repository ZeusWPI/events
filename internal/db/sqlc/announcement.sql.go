// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: announcement.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const announcementCreate = `-- name: AnnouncementCreate :one
INSERT INTO announcement (event_id, content, send_time, send, error)
VALUES ($1, $2, $3, $4, $5)
RETURNING id
`

type AnnouncementCreateParams struct {
	EventID  int32
	Content  string
	SendTime pgtype.Timestamptz
	Send     bool
	Error    pgtype.Text
}

func (q *Queries) AnnouncementCreate(ctx context.Context, arg AnnouncementCreateParams) (int32, error) {
	row := q.db.QueryRow(ctx, announcementCreate,
		arg.EventID,
		arg.Content,
		arg.SendTime,
		arg.Send,
		arg.Error,
	)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const announcementError = `-- name: AnnouncementError :exec
UPDATE announcement
SET error = $1
WHERE id = $2
`

type AnnouncementErrorParams struct {
	Error pgtype.Text
	ID    int32
}

func (q *Queries) AnnouncementError(ctx context.Context, arg AnnouncementErrorParams) error {
	_, err := q.db.Exec(ctx, announcementError, arg.Error, arg.ID)
	return err
}

const announcementGetByEvents = `-- name: AnnouncementGetByEvents :many
SELECT id, event_id, content, send_time, send, error 
FROM announcement
WHERE event_id = ANY($1::int[])
ORDER BY send_time
`

func (q *Queries) AnnouncementGetByEvents(ctx context.Context, dollar_1 []int32) ([]Announcement, error) {
	rows, err := q.db.Query(ctx, announcementGetByEvents, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Announcement
	for rows.Next() {
		var i Announcement
		if err := rows.Scan(
			&i.ID,
			&i.EventID,
			&i.Content,
			&i.SendTime,
			&i.Send,
			&i.Error,
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

const announcementGetUnsend = `-- name: AnnouncementGetUnsend :many
SELECT id, event_id, content, send_time, send, error
FROM announcement
WHERE NOT send AND error IS NULL
`

func (q *Queries) AnnouncementGetUnsend(ctx context.Context) ([]Announcement, error) {
	rows, err := q.db.Query(ctx, announcementGetUnsend)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Announcement
	for rows.Next() {
		var i Announcement
		if err := rows.Scan(
			&i.ID,
			&i.EventID,
			&i.Content,
			&i.SendTime,
			&i.Send,
			&i.Error,
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

const announcementSend = `-- name: AnnouncementSend :exec
UPDATE announcement
SET send = true
WHERE id = $1
`

func (q *Queries) AnnouncementSend(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, announcementSend, id)
	return err
}

const announcementUpdate = `-- name: AnnouncementUpdate :exec
UPDATE announcement
SET content = $1, send_time = $2
WHERE id = $3 AND NOT send
`

type AnnouncementUpdateParams struct {
	Content  string
	SendTime pgtype.Timestamptz
	ID       int32
}

func (q *Queries) AnnouncementUpdate(ctx context.Context, arg AnnouncementUpdateParams) error {
	_, err := q.db.Exec(ctx, announcementUpdate, arg.Content, arg.SendTime, arg.ID)
	return err
}
