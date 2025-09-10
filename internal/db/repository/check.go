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

// GetByUID returns the check with the given UID
// It searches inside the check table
// It does not fetch the associated check_events entries
func (c *Check) GetByUID(ctx context.Context, checkUID string) (*model.Check, error) {
	check, err := c.repo.queries(ctx).CheckGetByCheckUID(ctx, checkUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get check by uid %s | %w", checkUID, err)
	}

	return model.CheckModel(check, sqlc.CheckEvent{}), nil
}

// GetByID returns the check event with the given ID
// It looks inside the check_event table
// It also populate the associated 'check' table fields
func (c *Check) GetByID(ctx context.Context, checkID int) (*model.Check, error) {
	check, err := c.repo.queries(ctx).CheckGetByCheckEventID(ctx, int32(checkID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get check event by id %d | %w", checkID, err)
	}

	return model.CheckModel(check.Check, check.CheckEvent), nil
}

func (c *Check) GetByEvents(ctx context.Context, events []model.Event) ([]*model.Check, error) {
	checks, err := c.repo.queries(ctx).CheckGetByEvents(ctx, utils.SliceMap(events, func(e model.Event) int32 { return int32(e.ID) }))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get checks by events %+v | %w", events, err)
	}

	return utils.SliceMap(checks, func(check sqlc.CheckGetByEventsRow) *model.Check {
		return model.CheckModel(check.Check, check.CheckEvent)
	}), nil
}

func (c *Check) GetEventsByCheckUID(ctx context.Context, checkUID string) ([]*model.Check, error) {
	checks, err := c.repo.queries(ctx).CheckGetByCheckUIDAll(ctx, checkUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get check events by check id %s | %w", checkUID, err)
	}

	return utils.SliceMap(checks, func(check sqlc.CheckGetByCheckUIDAllRow) *model.Check {
		return model.CheckModel(check.Check, check.CheckEvent)
	}), nil
}

func (c *Check) GetByCheckUIDEvent(ctx context.Context, checkUID string, eventID int) (*model.Check, error) {
	check, err := c.repo.queries(ctx).CheckGetByCheckEvent(ctx, sqlc.CheckGetByCheckEventParams{
		Uid:     checkUID,
		EventID: int32(eventID),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get check by check id and event %s | %d | %w", checkUID, eventID, err)
	}

	return model.CheckModel(check.Check, check.CheckEvent), nil
}

// Create only creates the check, not the checkEvents
func (c *Check) Create(ctx context.Context, check model.Check) error {
	if err := c.repo.queries(ctx).CheckCreate(ctx, sqlc.CheckCreateParams{
		Uid:         check.UID,
		Description: check.Description,
		Deadline:    pgtype.Int8{Int64: check.Deadline.Nanoseconds(), Valid: check.Deadline.Nanoseconds() != 0},
		Active:      check.Active,
		Type:        sqlc.CheckType(check.Type),
		CreatorID:   pgtype.Int4{Int32: int32(check.CreatorID), Valid: check.CreatorID != 0},
	}); err != nil {
		return fmt.Errorf("create check %+v | %w", check, err)
	}

	return nil
}

func (c *Check) CreateEvent(ctx context.Context, check *model.Check) error {
	id, err := c.repo.queries(ctx).CheckEventCreate(ctx, sqlc.CheckEventCreateParams{
		CheckUid: check.UID,
		EventID:  int32(check.EventID),
		Status:   sqlc.CheckStatus(check.Status),
		Message:  pgtype.Text{String: check.Message, Valid: check.Message != ""},
	})
	if err != nil {
		return err
	}

	check.ID = int(id)

	return nil
}

// CreateEventBatch creates the check events in batch
func (c *Check) CreateEventBatch(ctx context.Context, checks []model.Check) error {
	if len(checks) == 0 {
		return nil
	}
	if checks[0].UID == "" {
		return fmt.Errorf("invalid check %+v", checks)
	}

	if err := c.repo.queries(ctx).CheckEventCreateBatch(ctx, sqlc.CheckEventCreateBatchParams{
		Column1: utils.SliceRepeat(checks[0].UID, len(checks)),
		Column2: utils.SliceMap(checks, func(check model.Check) int32 { return int32(check.EventID) }),
		Column3: utils.SliceMap(checks, func(check model.Check) string { return string(check.Status) }),
		Column4: utils.SliceMap(checks, func(check model.Check) string { return check.Message }),
	}); err != nil {
		return fmt.Errorf("create check event batch %+v | %w", checks, err)
	}

	return nil
}

// Update updates a check, not the check event
func (c *Check) Update(ctx context.Context, check model.Check) error {
	if err := c.repo.queries(ctx).CheckUpdate(ctx, sqlc.CheckUpdateParams{
		Uid:         check.UID,
		Description: check.Description,
		Deadline:    pgtype.Int8{Int64: check.Deadline.Nanoseconds(), Valid: check.Deadline.Nanoseconds() != 0},
		Active:      check.Active,
		Type:        sqlc.CheckType(check.Type),
	}); err != nil {
		return fmt.Errorf("update check %+v | %w", check, err)
	}

	return nil
}

// UpdateEvent updates an event check, not the actual check
func (c *Check) UpdateEvent(ctx context.Context, check model.Check) error {
	if err := c.repo.queries(ctx).CheckEventUpdate(ctx, sqlc.CheckEventUpdateParams{
		ID:      int32(check.ID),
		Status:  sqlc.CheckStatus(check.Status),
		Message: pgtype.Text{String: check.Message, Valid: check.Message != ""},
	}); err != nil {
		return fmt.Errorf("update check event %+v | %w", check, err)
	}

	return nil
}

func (c *Check) SetInactiveAutomatic(ctx context.Context) error {
	if err := c.repo.queries(ctx).CheckSetInactiveAutomatic(ctx); err != nil {
		return fmt.Errorf("set automatic checks to inactive %w", err)
	}

	return nil
}

func (c *Check) Delete(ctx context.Context, checkUID string) error {
	if err := c.repo.queries(ctx).CheckDelete(ctx, checkUID); err != nil {
		return fmt.Errorf("delete check %s | %w", checkUID, err)
	}

	return nil
}
