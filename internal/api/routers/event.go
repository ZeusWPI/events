package api

import (
	"github.com/ZeusWPI/events/internal/api/dto"
	"github.com/ZeusWPI/events/internal/api/service"
	"github.com/gofiber/fiber/v2"
)

type Event struct {
	router fiber.Router

	event service.Event
}

func NewEvent(router fiber.Router, service service.Service) *Event {
	api := &Event{
		router: router.Group("/event"),
		event:  *service.NewEvent(),
	}

	api.createRoutes()

	return api
}

func (r *Event) createRoutes() {
	r.router.Get("/year/:id", r.getByYear)
	r.router.Post("/organizers", r.updateOrganizers)
	r.router.Post("/sync", r.sync)
}

func (r *Event) getByYear(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	events, err := r.event.GetByYear(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(events)
}

func (r *Event) updateOrganizers(c *fiber.Ctx) error {
	var events []dto.Event
	if err := c.BodyParser(&events); err != nil {
		return fiber.ErrBadRequest
	}

	for _, event := range events {
		if err := dto.Validate.Struct(event); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	if err := r.event.UpdateOrganizers(c.Context(), events); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (r *Event) sync(c *fiber.Ctx) error {
	if err := r.event.Sync(); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusAccepted)
}
