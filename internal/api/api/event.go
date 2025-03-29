// Package api provides all routes
package api

import (
	"github.com/ZeusWPI/events/internal/api/dto"
	"github.com/ZeusWPI/events/internal/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// EventRouter contains all api routes related to events
type EventRouter struct {
	router fiber.Router

	event service.Event
}

// NewEventRouter creates a new event router
func NewEventRouter(service service.Service, router fiber.Router) *EventRouter {
	api := &EventRouter{
		router: router.Group("/event"),
		event:  service.NewEvent(),
	}
	api.createRoutes()

	return api
}

func (r *EventRouter) createRoutes() {
	r.router.Get("/year/:id", r.getByYear)
	r.router.Post("/organizers", r.updateOrganizers)
	r.router.Post("/sync", r.sync)
}

func (r *EventRouter) getByYear(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	events, err := r.event.GetByYear(c.Context(), dto.Year{ID: id})
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(events)
}

func (r *EventRouter) updateOrganizers(c *fiber.Ctx) error {
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
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (r *EventRouter) sync(c *fiber.Ctx) error {
	if err := r.event.Sync(); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return c.SendStatus(fiber.StatusAccepted)
}
