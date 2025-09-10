package server

import (
	"fmt"

	"github.com/ZeusWPI/events/internal/server/api"
	v1 "github.com/ZeusWPI/events/internal/server/api/v1"
	"github.com/ZeusWPI/events/internal/server/middleware"
	"github.com/ZeusWPI/events/internal/server/service"
	"github.com/ZeusWPI/events/internal/server/webhook"
	"github.com/ZeusWPI/events/pkg/config"
	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/postgres/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shareed2k/goth_fiber"
	"go.uber.org/zap"
)

type Server struct {
	Addr string
	App  *fiber.App
}

const port = 4000

func NewServer(service *service.Service, pool *pgxpool.Pool) *Server {
	app := fiber.New(fiber.Config{
		BodyLimit:      1024 * 1024 * 1024,
		ReadBufferSize: 8096,
	})

	app.Use(fiberzap.New(fiberzap.Config{
		Logger: zap.L(),
	}))

	if config.IsDev() {
		app.Use(cors.New(cors.Config{
			AllowOrigins:     "http://localhost:3000",
			AllowHeaders:     "Origin, Content-Type, Accept, Access-Control-Allow-Origin",
			AllowCredentials: true,
		}))
	}

	goth_fiber.SessionStore = session.New(session.Config{
		KeyLookup:      fmt.Sprintf("cookie:%s_session_id", config.GetDefaultString("app.name", "events")),
		CookieHTTPOnly: true,
		Storage:        postgres.New(postgres.Config{DB: pool}),
		CookieSecure:   !config.IsDev(),
	})

	// Swagger
	app.Use(swagger.New(swagger.Config{
		BasePath: "/api/v1",
		FilePath: "./docs/v1/swagger.json",
		CacheAge: 0,
	}))

	// Initialize all routes
	apiRouter := app.Group("/api")

	// Public api
	v1Router := apiRouter.Group("/v1")

	v1.NewEvent(v1Router, service)

	// Internal api
	api.NewAuth(apiRouter, service)

	protectedRouter := apiRouter.Use(middleware.ProtectedRoute)

	api.NewEvent(protectedRouter, service)
	api.NewYear(protectedRouter, service)
	api.NewOrganizer(protectedRouter, service)
	api.NewTask(protectedRouter, service)
	api.NewCheck(protectedRouter, service)
	api.NewAnnouncement(protectedRouter, service)
	api.NewMail(protectedRouter, service)
	api.NewPoster(protectedRouter, service)

	// Webhook
	webhookRouter := app.Group("/webhook")

	webhook.NewGithub(webhookRouter, service)
	webhook.NewGitmate(webhookRouter, service)

	if !config.IsDev() {
		app.Static("/", "./public")
		app.Static("*", "./public/index.html")
	}

	host := config.GetDefaultString("app.host", "0.0.0.0")

	return &Server{
		Addr: fmt.Sprintf("%s:%d", host, port),
		App:  app,
	}
}
