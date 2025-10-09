package v1

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

type imageID struct {
	ID int `json:"id"`
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
	r.router.Get("/:imageId", r.get)
	r.router.Put("/", r.create)
}

// get returns the image associated with an id
//
//	@Summary		Get an image
//	@Description	Get an image given an id
//	@Tags			image
//	@Produce		png
//	@Success		200
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Param			id	path	int	true	"image id""
//	@Router			/image/{id} [get]
func (r *Image) get(c *fiber.Ctx) error {
	id, err := c.ParamsInt("imageId")
	if err != nil {
		return fiber.ErrBadRequest
	}

	file, err := r.image.Get(c.Context(), id)
	if err != nil {
		return err
	}

	c.Set("Content-Type", mimePNG)

	return utils.SendCached(c, file)
}

// create creates a new image
//
//	@Summary		Store an image
//	@Description	Store the image and get the id back
//	@Tags			image
//	@Accept			x-www-form-urlencoded
//	@Produce		json
//	@Success		200
//	@Failure		400
//	@Failure		500
//	@Router			/image [put]
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

	return c.JSON(imageID{ID: id})
}
