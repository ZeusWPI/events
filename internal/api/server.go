// Package api provides all API routes and handlers
package api

import (
	"fmt"
	"strings"

	"github.com/ZeusWPI/events/internal/api/api"
	"github.com/ZeusWPI/events/internal/api/middleware"
	"github.com/ZeusWPI/events/internal/service"
	"github.com/ZeusWPI/events/pkg/config"
	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/postgres/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shareed2k/goth_fiber"
	"go.uber.org/zap"
)

// Server represents the api
type Server struct {
	Addr string
	App  *fiber.App
}

const port = 4000

// NewServer creates a new Server
func NewServer(service service.Service, pool *pgxpool.Pool) *Server {
	app := fiber.New(fiber.Config{
		BodyLimit:      1024 * 1024 * 1024,
		ReadBufferSize: 8096,
	})

	app.Use(fiberzap.New(fiberzap.Config{
		Logger: zap.L(),
	}))

	env := config.GetDefaultString("app.env", "development")
	env = strings.ToLower(env)

	if env == "development" {
		app.Use(cors.New(cors.Config{
			AllowOrigins:     "http://localhost:5173",
			AllowHeaders:     "Origin, Content-Type, Accept, Access-Control-Allow-Origin",
			AllowCredentials: true,
		}))
	}

	goth_fiber.SessionStore = session.New(session.Config{
		KeyLookup:      fmt.Sprintf("cookie:%s_session_id", config.GetDefaultString("app.name", "events")),
		CookieHTTPOnly: true,
		Storage:        postgres.New(postgres.Config{DB: pool}),
		CookieSecure:   env != "development",
	})

	// Initialize all routes
	apiRouter := app.Group("/api")
	api.NewAuthRouter(service, apiRouter)

	protectedRouter := apiRouter.Use(middleware.ProtectedRoute)

	api.NewEventRouter(service, protectedRouter)
	api.NewYearRouter(service, protectedRouter)
	api.NewOrganizerRouter(service, protectedRouter)
	api.NewTaskRouter(service, protectedRouter)

	if env != "development" {
		app.Static("/", "./public")
		app.Static("*", "./public/index.html")
	}

	host := config.GetDefaultString("app.host", "0.0.0.0")

	return &Server{
		Addr: fmt.Sprintf("%s:%d", host, port),
		App:  app,
	}
}
