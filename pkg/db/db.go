// Package db allows for the creation of a new postgres database connection
package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/ZeusWPI/events/internal/pkg/db/sqlc"
	"github.com/ZeusWPI/events/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Options are the different configuration options for a postgres connection
type Options struct {
	Host     string // Default: localhost
	Port     uint16 // Default: 5432
	Database string // Default: postgres
	User     string // Default: postgres
	Password string // Default: postgres
}

// DB represents a postgres database connection
type DB struct {
	Pool    *pgxpool.Pool
	Queries *sqlc.Queries
}

// New creates a new postgres database connection
func New(options Options) (*DB, error) {
	pgConfig, err := pgxpool.ParseConfig("")
	if err != nil {
		return nil, err
	}

	pgConfig.ConnConfig.Host = defaultString(options.Host, "localhost")
	pgConfig.ConnConfig.Port = defaultInt(options.Port, 5432)
	pgConfig.ConnConfig.Database = defaultString(options.Database, "postgres")
	pgConfig.ConnConfig.User = defaultString(options.User, "postgres")
	pgConfig.ConnConfig.Password = config.GetDefaultString(options.Password, "postgres")

	pool, err := pgxpool.NewWithConfig(context.Background(), pgConfig)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.TODO()); err != nil {
		return nil, err
	}

	queries := sqlc.New(pool)

	return &DB{Pool: pool, Queries: queries}, nil
}

// WithRollback executes a given operations function. If the function returns an error, all transactions are reverted
func (db *DB) WithRollback(ctx context.Context, operations func(context.Context, *sqlc.Queries) (interface{}, error)) (interface{}, error) {
	// Begin transaction
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("DB: Failed to begin transaction: %w", err)
	}

	q := db.Queries.WithTx(tx)

	// Execute the operations
	result, err := operations(ctx, q)
	if err != nil {
		// Rollback time!
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			err = errors.Join(err, fmt.Errorf("DB: Rollback failed: %w", rollbackErr))
		}

		return nil, fmt.Errorf("DB: Transactions failed: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("DB: Failed to commit transaction: %w", err)
	}

	return result, nil
}

func defaultString(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func defaultInt(value, defaultValue uint16) uint16 {
	if value == 0 {
		return defaultValue
	}
	return value
}
