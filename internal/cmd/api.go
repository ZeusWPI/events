package cmd

import (
	"fmt"

	"github.com/ZeusWPI/events/pkg/config"
	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
)

const port = 4000

// API starts the webserver serving the API and static files
func API() {
	app := fiber.New(fiber.Config{
		BodyLimit: 1024 * 1024 * 1024,
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

	host := config.GetDefaultString("app.host", "0.0.0.0")

	zap.S().Fatalf("Fatal server error %v", app.Listen(fmt.Sprintf("%s:%d", host, port)))
}
