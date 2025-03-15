// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: member.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const memberCreate = `-- name: MemberCreate :one
INSERT INTO member (name, username, zauth_id)
VALUES ($1, $2, $3)
RETURNING id
`

type MemberCreateParams struct {
	Name     string
	Username pgtype.Text
	ZauthID  pgtype.Int4
}

func (q *Queries) MemberCreate(ctx context.Context, arg MemberCreateParams) (int32, error) {
	row := q.db.QueryRow(ctx, memberCreate, arg.Name, arg.Username, arg.ZauthID)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const memberGetAll = `-- name: MemberGetAll :many
SELECT id, name, username, zauth_id FROM member
`

func (q *Queries) MemberGetAll(ctx context.Context) ([]Member, error) {
	rows, err := q.db.Query(ctx, memberGetAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Member
	for rows.Next() {
		var i Member
		if err := rows.Scan(
			&i.ID,
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

const memberGetByID = `-- name: MemberGetByID :one
SELECT id, name, username, zauth_id FROM member 
WHERE id = $1
`

func (q *Queries) MemberGetByID(ctx context.Context, id int32) (Member, error) {
	row := q.db.QueryRow(ctx, memberGetByID, id)
	var i Member
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Username,
		&i.ZauthID,
	)
	return i, err
}

const memberGetByName = `-- name: MemberGetByName :one
SELECT id, name, username, zauth_id FROM member 
WHERE name ILIKE $1
`

func (q *Queries) MemberGetByName(ctx context.Context, name string) (Member, error) {
	row := q.db.QueryRow(ctx, memberGetByName, name)
	var i Member
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Username,
		&i.ZauthID,
	)
	return i, err
}

const memberUpdate = `-- name: MemberUpdate :exec
UPDATE member 
SET name = $1, username = $2, zauth_id = $3
WHERE id = $4
`

type MemberUpdateParams struct {
	Name     string
	Username pgtype.Text
	ZauthID  pgtype.Int4
	ID       int32
}

func (q *Queries) MemberUpdate(ctx context.Context, arg MemberUpdateParams) error {
	_, err := q.db.Exec(ctx, memberUpdate,
		arg.Name,
		arg.Username,
		arg.ZauthID,
		arg.ID,
	)
	return err
}
