package api

import (
	"github.com/ZeusWPI/events/internal/api/dto"
	"github.com/ZeusWPI/events/internal/api/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Organizer contains all api routes related to organizers
type Organizer struct {
	router fiber.Router

	organizer service.Organizer
}

// NewOrganizer creates a new organizer router
func NewOrganizer(service service.Service, router fiber.Router) *Organizer {
	api := &Organizer{
		router:    router.Group("/organizer"),
		organizer: service.NewOrganizer(),
	}

	api.createRoutes()

	return api
}

func (r *Organizer) createRoutes() {
	r.router.Get("/year/:id", r.getByYear)
	r.router.Get("/me", r.meHandler)
}

func (r *Organizer) getByYear(c *fiber.Ctx) error {
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

func (r *Organizer) meHandler(c *fiber.Ctx) error {
	memberID := c.Locals("memberID").(int)

	user, err := r.organizer.GetByID(c.Context(), memberID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(user)
}
