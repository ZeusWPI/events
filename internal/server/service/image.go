package service

import (
	"bytes"
	"context"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/server/dto"
	"github.com/ZeusWPI/events/pkg/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Image struct {
	service Service

	image repository.Image
}

func (s *Service) NewImage() *Image {
	return &Image{
		service: *s,
		image:   *s.repo.NewImage(),
	}
}

func (i *Image) Get(ctx context.Context, imageID int) ([]byte, error) {
	image, err := i.image.Get(ctx, imageID)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}
	if image == nil {
		return nil, fiber.ErrNotFound
	}

	file, err := storage.S.Get(image.FileID)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}

	return file, nil
}

func (i *Image) Save(ctx context.Context, imageSave dto.ImageSave) (int, error) {
	if !isPng(imageSave.File) {
		return 0, fiber.NewError(fiber.StatusBadRequest, "Only PNG is supported")
	}

	image := model.Image{
		Name:   imageSave.Name,
		FileID: uuid.NewString(),
	}

	if err := i.service.withRollback(ctx, func(ctx context.Context) error {
		if err := i.image.Create(ctx, &image); err != nil {
			zap.S().Error(err)
			return fiber.ErrInternalServerError
		}

		if err := storage.S.Set(image.FileID, imageSave.File, 0); err != nil {
			zap.S().Error(err)
			return fiber.ErrInternalServerError
		}

		return nil
	}); err != nil {
		return 0, err
	}

	return image.ID, nil
}

func isPng(file []byte) bool {
	pngSignature := []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1A, '\n'}
	return len(file) >= len(pngSignature) && bytes.Equal(file[:len(pngSignature)], pngSignature)
}
