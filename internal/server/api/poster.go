package api

import (
	"github.com/ZeusWPI/events/internal/server/dto"
	"github.com/ZeusWPI/events/internal/server/service"
	"github.com/gofiber/fiber/v2"
)

type Poster struct {
	router fiber.Router

	poster service.Poster
}

func NewPoster(router fiber.Router, service *service.Service) *Poster {
	api := &Poster{
		router: router.Group("/poster"),
		poster: *service.NewPoster(),
	}

	api.createRoutes()

	return api
}

func (r *Poster) createRoutes() {
	r.router.Get("/:id/file", r.getFile)
	r.router.Delete("/:id", r.delete)
	r.router.Post("/:id", r.update)
	r.router.Put("/", r.create)
}

func (r *Poster) getFile(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	file, err := r.poster.GetFile(c.Context(), id)
	if err != nil {
		return err
	}

	c.Set("Content-Type", mimePNG)
	return c.Send(file)
}

func (r *Poster) create(c *fiber.Ctx) error {
	var poster dto.PosterSave
	if err := c.BodyParser(&poster); err != nil {
		return fiber.ErrBadRequest
	}
	if poster.ID != 0 {
		return fiber.ErrBadRequest
	}

	form, err := c.MultipartForm()
	if err != nil {
		return fiber.ErrBadRequest
	}
	file, err := getFormFile(form, "file")
	if err != nil {
		return err
	}

	poster.File = file

	if err := dto.Validate.Struct(poster); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if _, err := r.poster.Save(c.Context(), poster); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (r *Poster) update(c *fiber.Ctx) error {
	var poster dto.PosterSave
	if err := c.BodyParser(&poster); err != nil {
		return fiber.ErrBadRequest
	}
	if poster.ID == 0 {
		return fiber.ErrBadRequest
	}

	form, err := c.MultipartForm()
	if err != nil {
		return fiber.ErrBadRequest
	}
	file, err := getFormFile(form, "file")
	if err != nil {
		return err
	}

	poster.File = file

	if err := dto.Validate.Struct(poster); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if _, err := r.poster.Save(c.Context(), poster); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (r *Poster) delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	if err := r.poster.Delete(c.Context(), id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
