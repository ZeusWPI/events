// Entrypoint for the events application
package main

import (
	"github.com/ZeusWPI/events/internal/pkg/db/repository"
	"github.com/ZeusWPI/events/internal/pkg/website"
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
		zap.S().Fatal("Unable to connect to database", err)
	}

	// Temporarily populate db with all events and academic years
	repo := repository.New(db)
	w := website.New(*repo)
	err = w.UpdateAllAcademicYears()
	if err != nil {
		zap.S().Fatal("AcademicYear error", err)
	}
	err = w.UpdateAllEvents()
	if err != nil {
		zap.S().Fatal("Update error ", err)
	}
}
