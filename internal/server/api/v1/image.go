package v1

import (
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
	r.router.Get("/:imageId", r.get)
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
