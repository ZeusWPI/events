package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/sqlc"
	"github.com/ZeusWPI/events/pkg/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

type Check struct {
	repo Repository
}

func (r *Repository) NewCheck() *Check {
	return &Check{
		repo: *r,
	}
}

func (c *Check) GetCustom(ctx context.Context, checkID int) (*model.CheckCustom, error) {
	check, err := c.repo.queries(ctx).CheckCustomGet(ctx, int32(checkID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get custom check by id %d | %w", checkID, err)
	}

	return model.CheckCustomModel(check), nil
}

func (c *Check) GetAll(ctx context.Context) ([]*model.Check, error) {
	checks, err := c.repo.queries(ctx).CheckGetAll(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get all checks %w", err)
	}

	return utils.SliceMap(checks, model.CheckModel), nil
}

func (c *Check) GetStatusByEvent(ctx context.Context, eventID int) ([]*model.CheckStatus, error) {
	statusses, err := c.repo.queries(ctx).CheckStatusGetByEvent(ctx, int32(eventID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get all check statusses by event %d | %w", eventID, err)
	}

	return utils.SliceMap(statusses, model.CheckStatusModel), nil
}

func (c *Check) GetCustomByEvent(ctx context.Context, eventID int) ([]*model.CheckCustom, error) {
	customs, err := c.repo.queries(ctx).CheckCustomGetByEvent(ctx, int32(eventID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get all custom checks by event %d | %w", eventID, err)
	}

	return utils.SliceMap(customs, model.CheckCustomModel), nil
}

func (c *Check) Create(ctx context.Context, check *model.Check) error {
	id, err := c.repo.queries(ctx).CheckCreate(ctx, sqlc.CheckCreateParams{
		Description: check.Description,
		Deadline:    check.Deadline.Nanoseconds(),
	})
	if err != nil {
		return fmt.Errorf("create check %+v | %w", *check, err)
	}

	check.ID = int(id)

	return nil
}

func (c *Check) CreateStatus(ctx context.Context, status *model.CheckStatus) error {
	id, err := c.repo.queries(ctx).CheckStatusCreate(ctx, sqlc.CheckStatusCreateParams{
		EventID: int32(status.EventID),
		CheckID: int32(status.CheckID),
		Status:  sqlc.CheckStatusEnum(status.Status),
		Message: pgtype.Text{String: status.Message, Valid: status.Message != ""},
	})
	if err != nil {
		return fmt.Errorf("create check status %+v | %w", *status, err)
	}

	status.ID = int(id)

	return nil
}

func (c *Check) CreateCustom(ctx context.Context, custom *model.CheckCustom) error {
	id, err := c.repo.queries(ctx).CheckCustomCreate(ctx, sqlc.CheckCustomCreateParams{
		EventID:     int32(custom.EventID),
		Description: custom.Description,
		Status:      sqlc.CheckStatusEnum(custom.Status),
		CreatorID:   int32(custom.CreatorID),
	})
	if err != nil {
		return fmt.Errorf("create custom check %+v | %w", *custom, err)
	}

	custom.ID = int(id)

	return nil
}

func (c *Check) UpdateStatus(ctx context.Context, status model.CheckStatus) error {
	if err := c.repo.queries(ctx).CheckStatusUpdate(ctx, sqlc.CheckStatusUpdateParams{
		ID:      int32(status.ID),
		Status:  sqlc.CheckStatusEnum(status.Status),
		Message: pgtype.Text{String: status.Message, Valid: status.Message != ""},
	}); err != nil {
		return fmt.Errorf("update check status %+v | %w", status, err)
	}

	return nil
}

func (c *Check) UpdateCustom(ctx context.Context, custom model.CheckCustom) error {
	if err := c.repo.queries(ctx).CheckCustomUpdate(ctx, sqlc.CheckCustomUpdateParams{
		Description: custom.Description,
		Status:      sqlc.CheckStatusEnum(custom.Status),
	}); err != nil {
		return fmt.Errorf("update custom check %+v | %w", custom, err)
	}

	return nil
}
