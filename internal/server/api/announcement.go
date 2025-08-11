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

func NewAnnouncement(router fiber.Router, service *service.Service) *Announcement {
	api := &Announcement{
		router:       router.Group("/announcement"),
		announcement: *service.NewAnnouncement(),
	}

	api.createRoutes()

	return api
}

func (r *Announcement) createRoutes() {
	r.router.Get("/year/:id", r.getByYear)
	r.router.Put("/", r.create)
	r.router.Post("/:id", r.update)
}

func (r *Announcement) getByYear(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	announcements, err := r.announcement.GetByYear(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(announcements)
}

func (r *Announcement) create(c *fiber.Ctx) error {
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

func (r *Announcement) update(c *fiber.Ctx) error {
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
