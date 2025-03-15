// Entrypoint for the events application
package main

import (
	"github.com/ZeusWPI/events/internal/cmd"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/pkg/website"
	"github.com/ZeusWPI/events/internal/service"
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

	// Start website
	website := website.New(*repo)
	cmd.RunWebsitePeriodic(website)

	service := service.New(*repo)
	if err := cmd.API(*service, db.Pool()); err != nil {
		zap.S().Error(err)
	}
}
