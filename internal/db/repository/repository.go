// Package repository provides all repositories
package repository

import (
	"context"

	"github.com/ZeusWPI/events/internal/db/sqlc"
	"github.com/ZeusWPI/events/pkg/db"
)

// Repository is used to create specific repositories
type Repository struct {
	db db.DB
}

// Key used to store the queries object in the context
type contextKey string

const queryKey = contextKey("queries")

// New creates a new Repository
func New(db db.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) queries(ctx context.Context) *sqlc.Queries {
	if q, ok := ctx.Value(queryKey).(*sqlc.Queries); ok {
		return q
	}

	return r.db.Queries()
}

// WithRollback allows performing sequential database operations and rollbacks if one fails
func (r *Repository) WithRollback(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := ctx.Value(queryKey).(*sqlc.Queries); ok {
		// We're already in a rollback
		return fn(ctx)
	}

	return r.db.WithRollback(ctx, func(q *sqlc.Queries) error {
		txCtx := context.WithValue(ctx, queryKey, q)
		return fn(txCtx)
	})
}

// Table specific repositories

// NewYear creates a new Year repository
func (r *Repository) NewYear() Year {
	return &yearRepo{repo: *r}
}

// NewEvent creates a new Event repository
func (r *Repository) NewEvent() Event {
	return &eventRepo{repo: *r, organizer: r.NewOrganizer()}
}

// NewBoard creates a new Board repository
func (r *Repository) NewBoard() Board {
	return &boardRepo{repo: *r, member: r.NewMember(), year: r.NewYear()}
}

// NewMember creates a new Member repository
func (r *Repository) NewMember() Member {
	return &memberRepo{repo: *r}
}

// NewOrganizer creates a new Organizer repository
func (r *Repository) NewOrganizer() Organizer {
	return &organizerRepo{repo: *r}
}
