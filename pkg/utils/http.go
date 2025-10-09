package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func GetFormFile(form *multipart.Form, field string) ([]byte, error) {
	files := form.File[field]
	if len(files) == 0 {
		return nil, fiber.ErrBadRequest
	}

	file, err := files[0].Open()
	if err != nil {
		zap.S().Errorf("Failed to open file %v", err)
		return nil, fiber.ErrInternalServerError
	}

	content, err := io.ReadAll(file)
	if err != nil {
		zap.S().Errorf("Failed to read file %v", err)
		return nil, fiber.ErrInternalServerError
	}

	return content, nil
}

func SendCached(c *fiber.Ctx, img []byte) error {
	etag := fmt.Sprintf(`"%x"`, md5.Sum(img))

	if match := c.Get("If-None-Match"); match != "" {
		if match == etag {
			c.Status(fiber.StatusNotModified)
			return nil
		}
	}

	c.Set("Cache-Control", "public, max-age=86400, stale-while-revalidate=120")
	c.Set("ETag", etag)

	return c.Send(img)
}
