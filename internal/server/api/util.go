package api

import (
	"io"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	mimePNG = "image/png"
)

func getFormFile(form *multipart.Form, field string) ([]byte, error) {
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
