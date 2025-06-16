package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/sqlc"
	"github.com/ZeusWPI/events/pkg/utils"
)

type Check struct {
	repo Repository
}

func (r *Repository) NewCheck() *Check {
	return &Check{
		repo: *r,
	}
}

func (c *Check) GetByEvents(ctx context.Context, events []model.Event) ([]*model.Check, error) {
	checks, err := c.repo.queries(ctx).CheckGetByEvents(ctx, utils.SliceMap(events, func(e model.Event) int32 { return int32(e.ID) }))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get checks by events %+v | %w", events, err)
	}

	return utils.SliceMap(checks, model.CheckModel), nil
}

func (c *Check) Create(ctx context.Context, check *model.Check) error {
	id, err := c.repo.queries(ctx).CheckCreate(ctx, sqlc.CheckCreateParams{
		EventID:     int32(check.EventID),
		Description: check.Description,
		Done:        check.Done,
	})
	if err != nil {
		return fmt.Errorf("create check %+v | %w", *check, err)
	}

	check.ID = int(id)

	return nil
}

func (c *Check) Toggle(ctx context.Context, checkID int) error {
	if err := c.repo.queries(ctx).CheckToggle(ctx, int32(checkID)); err != nil {
		return fmt.Errorf("toggle check %d | %w", checkID, err)
	}

	return nil
}
