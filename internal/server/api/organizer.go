package api

import (
	"errors"

	"github.com/ZeusWPI/events/internal/server/service"
	"github.com/gofiber/fiber/v2"
)

type Organizer struct {
	router fiber.Router

	organizer service.Organizer
}

func NewOrganizer(router fiber.Router, service *service.Service) *Organizer {
	api := &Organizer{
		router:    router.Group("/organizer"),
		organizer: *service.NewOrganizer(),
	}

	api.createRoutes()

	return api
}

func (r *Organizer) createRoutes() {
	r.router.Get("/year/:id", r.getByYear)
	r.router.Get("/me", r.me)
}

func (r *Organizer) getByYear(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	organizers, err := r.organizer.GetByYear(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(organizers)
}

func (r *Organizer) me(c *fiber.Ctx) error {
	memberID, ok := c.Locals("memberID").(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	user, err := r.organizer.GetByMember(c.Context(), memberID)
	if err != nil {
		if errors.Is(err, fiber.ErrBadRequest) {
			return fiber.ErrForbidden
		}
		return err
	}

	return c.JSON(user)
}
