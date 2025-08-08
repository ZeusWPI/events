package webhook

import (
	"github.com/ZeusWPI/events/internal/server/service"
	"github.com/gofiber/fiber/v2"
)

type Gitmate struct {
	router fiber.Router
	poster service.Poster
}

func NewGitmate(router fiber.Router, service *service.Service) *Gitmate {
	webhook := &Gitmate{
		router: router.Group("/gitmate"),
		poster: *service.NewPoster(),
	}

	webhook.createRoutes()

	return webhook
}

func (r *Gitmate) createRoutes() {
	r.router.Post("/push", r.push)
}

func (r *Gitmate) push(c *fiber.Ctx) error {
	_ = r.poster.Sync()

	return c.SendStatus(fiber.StatusOK)
}
