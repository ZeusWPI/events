// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package sqlc

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Board struct {
	ID        int32
	Member    int32
	Year      int32
	Role      string
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

type Event struct {
	ID          int32
	Url         string
	Name        string
	Description pgtype.Text
	StartTime   pgtype.Timestamptz
	EndTime     pgtype.Timestamptz
	Location    pgtype.Text
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
	DeletedAt   pgtype.Timestamptz
	Year        int32
}

type Member struct {
	ID       int32
	Name     string
	Username pgtype.Text
	ZauthID  pgtype.Int4
}

type Organizer struct {
	ID    int32
	Event int32
	Board int32
}

type Task struct {
	ID        int32
	Name      string
	Result    pgtype.Text
	RunAt     pgtype.Timestamptz
	Error     pgtype.Text
	Recurring bool
}

type Year struct {
	ID        int32
	StartYear int32
	EndYear   int32
}
