// Entrypoint for the events application
package main

import (
	"github.com/ZeusWPI/events/internal/cmd"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/server/service"
	"github.com/ZeusWPI/events/internal/task"
	"github.com/ZeusWPI/events/internal/website"
	"github.com/ZeusWPI/events/pkg/config"
	"github.com/ZeusWPI/events/pkg/db"
	"github.com/ZeusWPI/events/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	err := config.Init()
	if err != nil {
		panic(err)
	}

	zapLogger := logger.New()
	zap.ReplaceGlobals(zapLogger)

	db, err := db.NewPSQL(db.PSQLOptions{
		Host:     config.GetString("db.host"),
		Port:     uint16(config.GetInt("db.port")),
		Database: config.GetString("db.database"),
		User:     config.GetString("db.user"),
		Password: config.GetString("db.password"),
	})
	if err != nil {
		zap.S().Fatalf("Unable to connect to database %v", err)
	}

	repo := repository.New(db)

	manager, err := task.NewManager(*repo)
	if err != nil {
		zap.S().Fatalf("Unable to create task manager %v", err)
	}

	// Start website
	website, err := website.New(*repo)
	if err != nil {
		zap.S().Fatalf("Unable to create website %v", err)
	}

	if err := cmd.Website(manager, *website); err != nil {
		zap.S().Fatalf("Unable to start website tasks %v", err)
	}

	// Start API
	service := service.New(*repo, manager, *website)
	if err := cmd.API(*service, db.Pool()); err != nil {
		zap.S().Error(err)
	}
}
