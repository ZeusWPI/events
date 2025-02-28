// Package api provides all routes
package api

import (
	"github.com/ZeusWPI/events/internal/api/dto"
	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/pkg/util"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// EventRouter contains all api routes related to events
type EventRouter struct {
	router fiber.Router

	event repository.Event
}

// NewEventRouter creates a new event router
func NewEventRouter(repo repository.Repository, router fiber.Router) *EventRouter {
	api := &EventRouter{
		router: router.Group("/event"),
		event:  repo.NewEvent(),
	}
	api.createRoutes()

	return api
}

func (r *EventRouter) createRoutes() {
	r.router.Get("/year/:id", r.getByYear)
}

func (r *EventRouter) getByYear(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	events, err := r.event.GetByYearWithAll(c.Context(), model.Year{ID: id})
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(util.SliceMap(events, dto.EventDTO))
}
