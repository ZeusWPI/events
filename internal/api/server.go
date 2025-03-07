// Package api provides all API routes and handlers
package api

import (
	"fmt"

	"github.com/ZeusWPI/events/internal/api/api"
	"github.com/ZeusWPI/events/internal/service"
	"github.com/ZeusWPI/events/pkg/config"
	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
)

// Server represents the api
type Server struct {
	Addr string
	App  *fiber.App
}

const port = 4000

// NewServer creates a new Server
func NewServer(service service.Service) *Server {
	app := fiber.New(fiber.Config{
		BodyLimit:      1024 * 1024 * 1024,
		ReadBufferSize: 8096,
	})
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: zap.L(),
	}))

	env := config.GetDefaultString("app.env", "development")
	if env != "development" {
		app.Static("/", "./public")
		app.Static("*", "./public/index.html")
	} else {
		app.Use(cors.New(cors.Config{
			AllowOrigins:     "http://localhost:5173",
			AllowHeaders:     "Origin, Content-Type, Accept, Access-Control-Allow-Origin",
			AllowCredentials: true,
		}))
	}

	apiRouter := app.Group("/api")

	// Initialize all routes
	api.NewEventRouter(service, apiRouter)
	api.NewYearRouter(service, apiRouter)
	api.NewOrganizerRouter(service, apiRouter)

	host := config.GetDefaultString("app.host", "0.0.0.0")
	return &Server{
		Addr: fmt.Sprintf("%s:%d", host, port),
		App:  app,
	}
}
