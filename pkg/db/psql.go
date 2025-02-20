package db

import (
	"context"
	"fmt"

	"github.com/ZeusWPI/events/internal/pkg/db/sqlc"
	"github.com/ZeusWPI/events/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PSQLOptions are the different configuration options for a postgres connection
type PSQLOptions struct {
	Host     string // Default: localhost
	Port     uint16 // Default: 5432
	Database string // Default: postgres
	User     string // Default: postgres
	Password string // Default: postgres
}

type psql struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

// Interface compliance
var _ DB = (*psql)(nil)

// NewPSQL creates a new postgres database connection
func NewPSQL(options PSQLOptions) (DB, error) {
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

	return &psql{pool: pool, queries: queries}, nil
}

// Queries returns a sqlc.Queries object
func (p *psql) Queries() *sqlc.Queries {
	return p.queries
}

// WithRollback perfoms the given function and rollbacks all performed transactions if one fails
func (p *psql) WithRollback(ctx context.Context, fn func(q *sqlc.Queries) error) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Failed to start transaction: %v", err)
	}
	defer func() {
		// Will error out if tx.Commit is called first
		_ = tx.Rollback(ctx)
	}()

	queries := sqlc.New(p.pool)

	if err := fn(queries); err != nil {
		return err
	}

	return tx.Commit(ctx)
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
