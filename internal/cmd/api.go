package cmd

import (
	"fmt"

	"github.com/ZeusWPI/events/internal/api"
	"github.com/ZeusWPI/events/internal/db/repository"
	"go.uber.org/zap"
)

// API starts the webserver serving the API and static files
func API(repo repository.Repository) error {
	server := api.NewServer(repo)

	zap.S().Infof("Server is running on %s", server.Addr)

	if err := server.App.Listen(server.Addr); err != nil {
		return fmt.Errorf("API unknown error %v", err)
	}

	return nil
}
