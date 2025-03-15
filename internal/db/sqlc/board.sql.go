// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: board.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const boardCreate = `-- name: BoardCreate :one
INSERT INTO board (member, year, role)
VALUES ($1, $2, $3)
RETURNING id
`

type BoardCreateParams struct {
	Member int32
	Year   int32
	Role   string
}

func (q *Queries) BoardCreate(ctx context.Context, arg BoardCreateParams) (int32, error) {
	row := q.db.QueryRow(ctx, boardCreate, arg.Member, arg.Year, arg.Role)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const boardGetAllWithMemberYear = `-- name: BoardGetAllWithMemberYear :many
SELECT b.id, member, year, role, created_at, updated_at, m.id, name, username, zauth_id, y.id, start_year, end_year FROM board b 
INNER JOIN member m ON b.member = m.id 
INNER JOIN year y ON b.year = y.id
`

type BoardGetAllWithMemberYearRow struct {
	ID        int32
	Member    int32
	Year      int32
	Role      string
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	ID_2      int32
	Name      string
	Username  pgtype.Text
	ZauthID   pgtype.Int4
	ID_3      int32
	StartYear int32
	EndYear   int32
}

func (q *Queries) BoardGetAllWithMemberYear(ctx context.Context) ([]BoardGetAllWithMemberYearRow, error) {
	rows, err := q.db.Query(ctx, boardGetAllWithMemberYear)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []BoardGetAllWithMemberYearRow
	for rows.Next() {
		var i BoardGetAllWithMemberYearRow
		if err := rows.Scan(
			&i.ID,
			&i.Member,
			&i.Year,
			&i.Role,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ID_2,
			&i.Name,
			&i.Username,
			&i.ZauthID,
			&i.ID_3,
			&i.StartYear,
			&i.EndYear,
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

const boardGetByMemberYear = `-- name: BoardGetByMemberYear :one
SELECT b.id, member, year, role, created_at, updated_at, m.id, name, username, zauth_id, y.id, start_year, end_year FROM board b 
INNER JOIN member m ON b.member = m.id 
INNER JOIN year y ON b.year = y.id
WHERE m.id = $1 AND y.id = $2
`

type BoardGetByMemberYearParams struct {
	ID   int32
	ID_2 int32
}

type BoardGetByMemberYearRow struct {
	ID        int32
	Member    int32
	Year      int32
	Role      string
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	ID_2      int32
	Name      string
	Username  pgtype.Text
	ZauthID   pgtype.Int4
	ID_3      int32
	StartYear int32
	EndYear   int32
}

func (q *Queries) BoardGetByMemberYear(ctx context.Context, arg BoardGetByMemberYearParams) (BoardGetByMemberYearRow, error) {
	row := q.db.QueryRow(ctx, boardGetByMemberYear, arg.ID, arg.ID_2)
	var i BoardGetByMemberYearRow
	err := row.Scan(
		&i.ID,
		&i.Member,
		&i.Year,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ID_2,
		&i.Name,
		&i.Username,
		&i.ZauthID,
		&i.ID_3,
		&i.StartYear,
		&i.EndYear,
	)
	return i, err
}

const boardGetByYearWithMemberYear = `-- name: BoardGetByYearWithMemberYear :many
SELECT b.id, member, year, role, created_at, updated_at, m.id, name, username, zauth_id, y.id, start_year, end_year FROM board b 
INNER JOIN member m ON b.member = m.id 
INNER JOIN year y ON b.year = y.id
WHERE b.year = $1
`

type BoardGetByYearWithMemberYearRow struct {
	ID        int32
	Member    int32
	Year      int32
	Role      string
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	ID_2      int32
	Name      string
	Username  pgtype.Text
	ZauthID   pgtype.Int4
	ID_3      int32
	StartYear int32
	EndYear   int32
}

func (q *Queries) BoardGetByYearWithMemberYear(ctx context.Context, year int32) ([]BoardGetByYearWithMemberYearRow, error) {
	rows, err := q.db.Query(ctx, boardGetByYearWithMemberYear, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []BoardGetByYearWithMemberYearRow
	for rows.Next() {
		var i BoardGetByYearWithMemberYearRow
		if err := rows.Scan(
			&i.ID,
			&i.Member,
			&i.Year,
			&i.Role,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ID_2,
			&i.Name,
			&i.Username,
			&i.ZauthID,
			&i.ID_3,
			&i.StartYear,
			&i.EndYear,
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
