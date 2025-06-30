// Entrypoint for the events application
package main

import (
	"github.com/ZeusWPI/events/internal/check"
	"github.com/ZeusWPI/events/internal/cmd"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/dsa"
	"github.com/ZeusWPI/events/internal/mail"
	"github.com/ZeusWPI/events/internal/mattermost"
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

	checkManager := check.NewManager(*repo)

	taskManager, err := task.NewManager(*repo)
	if err != nil {
		zap.S().Fatalf("Unable to create task manager %v", err)
	}

	// Start website
	website, err := website.New(*repo)
	if err != nil {
		zap.S().Fatalf("Unable to create website %v", err)
	}

	if err := cmd.Website(*website, taskManager); err != nil {
		zap.S().Fatalf("Unable to start website command %v", err)
	}

	// Start dsa
	dsa, err := dsa.New(*repo)
	if err != nil {
		zap.S().Fatalf("Unable to create dsa %v", err)
	}

	if err := cmd.DSA(dsa, taskManager, checkManager); err != nil {
		zap.S().Fatalf("Unable to start dsa command %v", err)
	}

	// Start mattermost
	mattermost, err := mattermost.New(*repo, taskManager)
	if err != nil {
		zap.S().Fatalf("Unable to create mattermost %v", err)
	}

	if err := cmd.Mattermost(mattermost, checkManager); err != nil {
		zap.S().Fatalf("Unable to start mattermost command %v", err)
	}

	// Start mail
	mail, err := mail.New(*repo, taskManager)
	if err != nil {
		zap.S().Fatalf("Unable to create mail %v", err)
	}

	if err := cmd.Mail(mail, checkManager); err != nil {
		zap.S().Fatalf("Unable to start mail command %v", err)
	}

	// Start API
	service := service.New(*repo, checkManager, taskManager, *mail, *website, *mattermost)
	if err := cmd.API(*service, db.Pool()); err != nil {
		zap.S().Error(err)
	}
}
