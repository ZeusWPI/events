package api

import (
	"github.com/ZeusWPI/events/internal/server/dto"
	"github.com/ZeusWPI/events/internal/server/service"
	"github.com/ZeusWPI/events/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type Image struct {
	router fiber.Router

	image service.Image
}

func NewImage(router fiber.Router, service *service.Service) *Image {
	api := &Image{
		router: router.Group("/image"),
		image:  *service.NewImage(),
	}

	api.createRoutes()

	return api
}

func (r *Image) createRoutes() {
	r.router.Put("/", r.create)
}

func (r *Image) create(c *fiber.Ctx) error {
	var image dto.ImageSave
	if err := c.BodyParser(&image); err != nil {
		return fiber.ErrBadRequest
	}

	form, err := c.MultipartForm()
	if err != nil {
		return fiber.ErrBadRequest
	}
	file, err := utils.GetFormFile(form, "file")
	if err != nil {
		return err
	}

	image.File = file

	if err := dto.Validate.Struct(image); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	id, err := r.image.Save(c.Context(), image)
	if err != nil {
		return err
	}

	return c.JSON(dto.ImageID{ID: id})
}
