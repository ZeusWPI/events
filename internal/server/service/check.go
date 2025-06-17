package service

import (
	"context"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/server/dto"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Check struct {
	service Service

	check repository.Check
}

func (s *Service) NewCheck() *Check {
	return &Check{
		service: *s,
		check:   *s.repo.NewCheck(),
	}
}

func (c *Check) Create(ctx context.Context, checkSave *dto.Check) error {
	check := model.Check{
		EventID:     checkSave.EventID,
		Description: checkSave.Description,
	}

	if err := c.check.Create(ctx, &check); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return nil
}

func (c *Check) Toggle(ctx context.Context, checkID int) error {
	if err := c.check.Toggle(ctx, checkID); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return nil
}

func (c *Check) Delete(ctx context.Context, checkID int) error {
	if err := c.check.Delete(ctx, checkID); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return nil
}
