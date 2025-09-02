package api

import (
	"github.com/ZeusWPI/events/internal/server/dto"
	"github.com/ZeusWPI/events/internal/server/service"
	"github.com/gofiber/fiber/v2"
)

type Mail struct {
	router fiber.Router

	mail service.Mail
}

func NewMail(router fiber.Router, service *service.Service) *Mail {
	api := &Mail{
		router: router.Group("/mail"),
		mail:   *service.NewMail(),
	}

	api.createRoutes()

	return api
}

func (r *Mail) createRoutes() {
	r.router.Get("/year/:id", r.getByYear)
	r.router.Put("/", r.create)
	r.router.Post("/resend/:id", r.resend)
	r.router.Post("/:id", r.update)
	r.router.Delete("/:id", r.delete)
}

func (r *Mail) getByYear(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	mails, err := r.mail.GetByYear(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(mails)
}

func (r *Mail) create(c *fiber.Ctx) error {
	var mail dto.Mail
	if err := c.BodyParser(&mail); err != nil {
		return fiber.ErrBadRequest
	}
	if mail.ID != 0 {
		return fiber.ErrBadRequest
	}

	if err := dto.Validate.Struct(mail); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	userID, ok := c.Locals("memberID").(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	if _, err := r.mail.Save(c.Context(), mail, userID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (r *Mail) update(c *fiber.Ctx) error {
	var mail dto.Mail
	if err := c.BodyParser(&mail); err != nil {
		return fiber.ErrBadRequest
	}
	if mail.ID == 0 {
		return fiber.ErrBadRequest
	}

	if err := dto.Validate.Struct(mail); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	userID, ok := c.Locals("memberID").(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	if _, err := r.mail.Save(c.Context(), mail, userID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func (r *Mail) delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	if err := r.mail.Delete(c.Context(), id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (r *Mail) resend(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	userID, ok := c.Locals("memberID").(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	if err := r.mail.Resend(c.Context(), id, userID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
