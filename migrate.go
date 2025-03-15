// Package main runs all pending migrations
package main

import (
	"embed"

	"github.com/ZeusWPI/events/pkg/config"
	"github.com/ZeusWPI/events/pkg/db"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed db/migrations/*.sql
var embedMigrations embed.FS

func main() {
	if err := config.Init(); err != nil {
		panic(err)
	}
	// setup database
	database, err := db.NewPSQL(db.PSQLOptions{
		Host:     config.GetString("db.host"),
		Port:     uint16(config.GetInt("db.port")),
		Database: config.GetString("db.database"),
		User:     config.GetString("db.user"),
		Password: config.GetString("db.password"),
	})
	if err != nil {
		panic(err)
	}
	db := stdlib.OpenDBFromPool(database.Pool())

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "db/migrations"); err != nil {
		panic(err)
	}

	// run app
}
