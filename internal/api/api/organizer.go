package api

import (
	"github.com/ZeusWPI/events/internal/api/dto"
	"github.com/ZeusWPI/events/internal/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// OrganizerRouter contains all api routes related to organizers
type OrganizerRouter struct {
	router fiber.Router

	organizer service.Organizer
}

// NewOrganizerRouter creates a new organizer router
func NewOrganizerRouter(service service.Service, router fiber.Router) *OrganizerRouter {
	api := &OrganizerRouter{
		router:    router.Group("/organizer"),
		organizer: service.NewOrganizer(),
	}

	api.createRoutes()

	return api
}

func (r *OrganizerRouter) createRoutes() {
	r.router.Get("/year/:id", r.getByYear)
}

func (r *OrganizerRouter) getByYear(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	organizers, err := r.organizer.GetByYear(c.Context(), dto.Year{ID: id})
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(organizers)
}
