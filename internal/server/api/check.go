package api

import (
	"github.com/ZeusWPI/events/internal/server/dto"
	"github.com/ZeusWPI/events/internal/server/service"
	"github.com/gofiber/fiber/v2"
)

type Check struct {
	router fiber.Router

	check service.Check
}

func NewCheck(router fiber.Router, service service.Service) *Check {
	api := &Check{
		router: router.Group("/check"),
		check:  *service.NewCheck(),
	}

	api.createRoutes()

	return api
}

func (r *Check) createRoutes() {
	r.router.Post("/", r.create)
	r.router.Post("/:id", r.toggle)
}

func (r *Check) create(c *fiber.Ctx) error {
	var check dto.Check
	if err := c.BodyParser(&check); err != nil {
		return fiber.ErrBadRequest
	}

	if err := dto.Validate.Struct(check); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := r.check.Create(c.Context(), &check); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (r *Check) toggle(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	if err := r.check.Toggle(c.Context(), id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
