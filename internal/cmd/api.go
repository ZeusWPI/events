package cmd

import (
	"fmt"

	"github.com/ZeusWPI/events/internal/api"
	"github.com/ZeusWPI/events/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// API starts the webserver serving the API and static files
func API(service service.Service, pool *pgxpool.Pool) error {
	server := api.NewServer(service, pool)

	zap.S().Infof("Server is running on %s", server.Addr)

	if err := server.App.Listen(server.Addr); err != nil {
		return fmt.Errorf("API unknown error %v", err)
	}

	return nil
}
