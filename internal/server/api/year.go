package api

import (
	"github.com/ZeusWPI/events/internal/server/service"
	"github.com/gofiber/fiber/v2"
)

type Year struct {
	router fiber.Router

	year service.Year
}

func NewYear(router fiber.Router, service *service.Service) *Year {
	api := &Year{
		router: router.Group("/year"),
		year:   *service.NewYear(),
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
		return err
	}

	return c.JSON(years)
}
