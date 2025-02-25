// Package db provides logic to connect to a database
package db

import (
	"context"

	"github.com/ZeusWPI/events/internal/db/sqlc"
)

// DB represents a database connection
type DB interface {
	Queries() *sqlc.Queries
	WithRollback(ctx context.Context, fn func(q *sqlc.Queries) error) error
}
