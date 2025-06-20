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

func NewMail(router fiber.Router, service service.Service) *Mail {
	api := &Mail{
		router: router.Group("/mail"),
		mail:   *service.NewMail(),
	}

	api.createRoutes()

	return api
}

func (r *Mail) createRoutes() {
	r.router.Get("/", r.GetAll)
	r.router.Put("/", r.Create)
	r.router.Post("/:id", r.Update)
}

func (r *Mail) GetAll(c *fiber.Ctx) error {
	mails, err := r.mail.GetAll(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(mails)
}

func (r *Mail) Create(c *fiber.Ctx) error {
	var mail dto.MailSave
	if err := c.BodyParser(&mail); err != nil {
		return fiber.ErrBadRequest
	}
	if mail.ID != 0 {
		return fiber.ErrBadRequest
	}

	if err := dto.Validate.Struct(mail); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if _, err := r.mail.Save(c.Context(), mail); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (r *Mail) Update(c *fiber.Ctx) error {
	var mail dto.MailSave
	if err := c.BodyParser(&mail); err != nil {
		return fiber.ErrBadRequest
	}
	if mail.ID == 0 {
		return fiber.ErrBadRequest
	}

	if err := dto.Validate.Struct(mail); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := dto.Validate.Struct(mail); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if _, err := r.mail.Save(c.Context(), mail); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
