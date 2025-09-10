// Entrypoint for the events application
package main

import (
	"github.com/ZeusWPI/events/internal/announcement"
	"github.com/ZeusWPI/events/internal/check"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/dsa"
	"github.com/ZeusWPI/events/internal/mail"
	"github.com/ZeusWPI/events/internal/poster"
	"github.com/ZeusWPI/events/internal/server"
	"github.com/ZeusWPI/events/internal/server/service"
	"github.com/ZeusWPI/events/internal/task"
	"github.com/ZeusWPI/events/internal/website"
	"github.com/ZeusWPI/events/pkg/config"
	"github.com/ZeusWPI/events/pkg/db"
	"github.com/ZeusWPI/events/pkg/logger"
	"github.com/ZeusWPI/events/pkg/storage"
	"go.uber.org/zap"
)

func main() {
	err := config.Init()
	if err != nil {
		panic(err)
	}

	zapLogger := logger.New()
	zap.ReplaceGlobals(zapLogger)

	// Databases
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

	if err = storage.Init(db.Pool()); err != nil {
		zap.S().Fatalf("Unable to init storage %v", err)
	}

	repo := repository.New(db)

	if err := task.Init(*repo); err != nil {
		zap.S().Fatalf("Unable to initialize the task manager %v", err)
	}

	if err := check.Init(*repo); err != nil {
		zap.S().Fatalf("Unable to initialize the check manager %v", err)
	}

	// Start dsa
	dsa, err := dsa.New(*repo)
	if err != nil {
		zap.S().Fatalf("Unable to create dsa %v", err)
	}

	// Start website
	website, err := website.New(*repo, *dsa)
	if err != nil {
		zap.S().Fatalf("Unable to create website %v", err)
	}

	// Start announcement
	announcement, err := announcement.New(*repo)
	if err != nil {
		zap.S().Fatalf("Unable to create mattermost %v", err)
	}

	// Start mail
	mail, err := mail.New(*repo)
	if err != nil {
		zap.S().Fatalf("Unable to create mail %v", err)
	}

	// Start poster
	poster, err := poster.New(*repo)
	if err != nil {
		zap.S().Fatalf("Unable to create poster %v", err)
	}

	// Start API
	service := service.New(*repo, *mail, website, *announcement, *poster)
	server := server.NewServer(service, db.Pool())

	if err := server.App.Listen(server.Addr); err != nil {
		zap.S().Errorf("API unknown error %w", err)
	}
}
