package webhook

import (
	"github.com/ZeusWPI/events/internal/server/service"
	"github.com/gofiber/fiber/v2"
)

type Github struct {
	router  fiber.Router
	website service.Website
}

func NewGithub(router fiber.Router, service service.Service) *Github {
	webhook := &Github{
		router:  router.Group("/github"),
		website: *service.NewWebsite(),
	}

	webhook.createRoutes()

	return webhook
}

func (r *Github) createRoutes() {
	r.router.Post("/push", r.push)
}

type pushPayload struct {
	Ref string `json:"ref"`
}

func (r *Github) push(c *fiber.Ctx) error {
	var payload pushPayload
	if err := c.BodyParser(&payload); err != nil {
		return fiber.ErrBadRequest
	}

	if payload.Ref != "refs/heads/master" {
		return c.SendStatus(fiber.StatusNoContent)
	}

	_ = r.website.Sync()

	return c.SendStatus(fiber.StatusOK)
}
