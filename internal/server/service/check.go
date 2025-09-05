package service

import (
	"context"

	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/server/dto"
	"github.com/ZeusWPI/events/pkg/utils"
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

func (c *Check) Create(ctx context.Context, create dto.CheckCreate) error {
	check := create.ToModel()
}

func (c *Check) Update(ctx context.Context, update dto.CheckUpdate) error {
	check := update.ToModel()

	original, err := c.check.GetCustom(ctx, check.ID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}
	if original == nil {
		return fiber.ErrNotFound
	}

	utils.Merge(check, *original)

	if err := c.check.UpdateCustom(ctx, *check); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return nil
}
