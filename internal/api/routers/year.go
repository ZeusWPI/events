package api

import (
	"github.com/ZeusWPI/events/internal/api/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Year contains all api routes related to year
type Year struct {
	router fiber.Router

	year service.Year
}

// NewYear constructs a new year router
func NewYear(service service.Service, router fiber.Router) *Year {
	api := &Year{
		router: router.Group("/year"),
		year:   service.NewYear(),
	}
	api.createRoutes()

	return api
}

func (r *Year) createRoutes() {
	r.router.Get("/", r.getAll)
}

func (r *Year) getAll(c *fiber.Ctx) error {
	years, err := r.year.GetAll(c.Context())
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(years)
}
