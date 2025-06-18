package api

import (
	"github.com/ZeusWPI/events/internal/server/dto"
	"github.com/ZeusWPI/events/internal/server/service"
	"github.com/gofiber/fiber/v2"
)

type Announcement struct {
	router fiber.Router

	announcement service.Announcement
}

func NewAnnouncement(router fiber.Router, service service.Service) *Announcement {
	api := &Announcement{
		router:       router.Group("/announcement"),
		announcement: *service.NewAnnouncement(),
	}

	api.createRoutes()

	return api
}

func (r *Announcement) createRoutes() {
	r.router.Put("/", r.Create)
	r.router.Post("/:id", r.Update)
}

func (r *Announcement) Create(c *fiber.Ctx) error {
	var announcement dto.Announcement
	if err := c.BodyParser(&announcement); err != nil {
		return fiber.ErrBadRequest
	}
	if announcement.ID != 0 {
		return fiber.ErrBadRequest
	}

	if err := dto.Validate.Struct(announcement); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if _, err := r.announcement.Save(c.Context(), announcement); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (r *Announcement) Update(c *fiber.Ctx) error {
	var announcement dto.Announcement
	if err := c.BodyParser(&announcement); err != nil {
		return fiber.ErrBadRequest
	}
	if announcement.ID == 0 {
		return fiber.ErrBadRequest
	}

	if err := dto.Validate.Struct(announcement); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if _, err := r.announcement.Save(c.Context(), announcement); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
