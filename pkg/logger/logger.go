// Package logger provides a logger instance
package logger

import (
	"strings"

	"github.com/ZeusWPI/events/pkg/config"
	"go.uber.org/zap"
)

// New returns a new logger instance
func New() *zap.Logger {
	var logger *zap.Logger
	env := config.GetDefaultString("app.env", "development")
	env = strings.ToLower(env)

	if env == "development" {
		logger = zap.Must(zap.NewDevelopment(zap.AddStacktrace(zap.WarnLevel)))
	} else {
		logger = zap.Must(zap.NewProduction())
	}
	logger = logger.With(zap.String("env", env))

	return logger
}
