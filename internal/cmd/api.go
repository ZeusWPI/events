package cmd

import (
	"fmt"

	"github.com/ZeusWPI/events/internal/server"
	"github.com/ZeusWPI/events/internal/server/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// API starts the webserver serving the API and static files
func API(service *service.Service, pool *pgxpool.Pool) error {
	server := server.NewServer(service, pool)

	zap.S().Infof("Server is running on %s", server.Addr)

	if err := server.App.Listen(server.Addr); err != nil {
		return fmt.Errorf("API unknown error %w", err)
	}

	return nil
}
