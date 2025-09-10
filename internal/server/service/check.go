package service

import (
	"context"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/server/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Check struct {
	service Service

	board repository.Board
	check repository.Check
	year  repository.Year
}

func (s *Service) NewCheck() *Check {
	return &Check{
		service: *s,
		check:   *s.repo.NewCheck(),
	}
}

func (c *Check) Create(ctx context.Context, checkSave dto.Check, memberID int) (dto.Check, error) {
	year, err := c.year.GetLast(ctx)
	if err != nil {
		zap.S().Error(err)
		return dto.Check{}, fiber.ErrInternalServerError
	}
	if year == nil {
		// No.
		return dto.Check{}, fiber.ErrInternalServerError
	}

	board, err := c.board.GetByMemberYear(ctx, memberID, year.ID)
	if err != nil {
		zap.S().Error(err)
		return dto.Check{}, fiber.ErrInternalServerError
	}
	if board == nil {
		return dto.Check{}, fiber.ErrForbidden
	}

	check := model.Check{
		UID:         uuid.NewString(),
		Description: checkSave.Description,
		Active:      true,
		Type:        model.CheckManual,
		CreatorID:   board.ID,
		EventID:     checkSave.EventID,
		Status:      checkSave.Status,
		Message:     checkSave.Message,
	}

	if err := c.service.withRollback(ctx, func(ctx context.Context) error {
		if err := c.check.Create(ctx, check); err != nil {
			zap.S().Error(err)
			return fiber.ErrInternalServerError
		}

		if err := c.check.CreateEvent(ctx, &check); err != nil {
			zap.S().Error(err)
			return fiber.ErrInternalServerError
		}

		return nil
	}); err != nil {
		return dto.Check{}, err
	}

	return dto.CheckDTO(&check), nil
}

func (c *Check) Update(ctx context.Context, checkSave dto.CheckUpdate) (dto.Check, error) {
	check, err := c.check.GetByID(ctx, checkSave.ID)
	if err != nil {
		zap.S().Error(err)
		return dto.Check{}, fiber.ErrInternalServerError
	}
	if check == nil {
		return dto.Check{}, fiber.ErrBadRequest
	}
	if check.Type != model.CheckManual {
		return dto.Check{}, fiber.ErrBadRequest
	}

	check.Status = checkSave.Status
	check.Message = checkSave.Message

	if err := c.check.UpdateEvent(ctx, *check); err != nil {
		zap.S().Error(err)
		return dto.Check{}, fiber.ErrInternalServerError
	}

	return dto.CheckDTO(check), nil
}

func (c *Check) Delete(ctx context.Context, checkID int) error {
	check, err := c.check.GetByID(ctx, checkID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}
	if check == nil {
		return fiber.ErrBadRequest
	}
	if check.Type != model.CheckManual {
		return fiber.ErrBadRequest
	}

	// This will also delete the associated check event by cascade
	if err := c.check.Delete(ctx, check.UID); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return nil
}
