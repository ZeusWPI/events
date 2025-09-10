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

type Event struct {
	repo Repository

	organizer Organizer
}

func (r *Repository) NewEvent() *Event {
	return &Event{
		repo:      *r,
		organizer: *r.NewOrganizer(),
	}
}

func (e *Event) GetByID(ctx context.Context, eventID int) (*model.Event, error) {
	event, err := e.repo.queries(ctx).EventGet(ctx, int32(eventID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get event by id %d | %w", eventID, err)
	}

	return model.EventModel(event.Event, event.Year), nil
}

func (e *Event) GetByIDs(ctx context.Context, eventIDs []int) ([]*model.Event, error) {
	events, err := e.repo.queries(ctx).EventGetByIds(ctx, utils.SliceMap(eventIDs, func(id int) int32 { return int32(id) }))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get event by ids %+v | %w", eventIDs, err)
	}

	return utils.SliceMap(events, func(event sqlc.EventGetByIdsRow) *model.Event {
		return model.EventModel(event.Event, event.Year)
	}), nil
}

func (e *Event) GetAll(ctx context.Context) ([]*model.Event, error) {
	events, err := e.repo.queries(ctx).EventGetAll(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get events %w", err)
	}

	return utils.SliceMap(events, func(event sqlc.EventGetAllRow) *model.Event {
		return model.EventModel(event.Event, event.Year)
	}), nil
}

func (e *Event) GetByYear(ctx context.Context, yearID int) ([]*model.Event, error) {
	events, err := e.repo.queries(ctx).EventGetByYear(ctx, int32(yearID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get events by year %d | %w", yearID, err)
	}

	return utils.SliceMap(events, func(event sqlc.EventGetByYearRow) *model.Event {
		return model.EventModel(event.Event, event.Year)
	}), nil
}

func (e *Event) GetFuture(ctx context.Context) ([]*model.Event, error) {
	events, err := e.repo.queries(ctx).EventGetFuture(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get future events %w", err)
	}

	return utils.SliceMap(events, func(event sqlc.EventGetFutureRow) *model.Event {
		return model.EventModel(event.Event, event.Year)
	}), nil
}

func (e *Event) GetNext(ctx context.Context) (*model.Event, error) {
	event, err := e.repo.queries(ctx).EventGetNext(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get next event %w", err)
	}

	return model.EventModel(event.Event, event.Year), nil
}

func (e *Event) Create(ctx context.Context, event *model.Event) error {
	id, err := e.repo.queries(ctx).EventCreate(ctx, sqlc.EventCreateParams{
		FileName:    event.FileName,
		Name:        event.Name,
		Description: pgtype.Text{String: event.Description, Valid: true},
		StartTime:   pgtype.Timestamptz{Time: event.StartTime, Valid: true},
		EndTime:     pgtype.Timestamptz{Time: event.EndTime, Valid: !event.EndTime.IsZero()},
		YearID:      int32(event.YearID),
		Location:    pgtype.Text{String: event.Location, Valid: true},
		Deleted:     event.Deleted,
	})
	if err != nil {
		return fmt.Errorf("create event %+v | %w", *event, err)
	}

	event.ID = int(id)

	return nil
}

func (e *Event) Update(ctx context.Context, event model.Event) error {
	if err := e.repo.queries(ctx).EventUpdate(ctx, sqlc.EventUpdateParams{
		ID:          int32(event.ID),
		Name:        event.Name,
		Description: pgtype.Text{String: event.Description, Valid: event.Description != ""},
		StartTime:   pgtype.Timestamptz{Time: event.StartTime, Valid: true},
		EndTime:     pgtype.Timestamptz{Time: event.EndTime, Valid: !event.EndTime.IsZero()},
		YearID:      int32(event.YearID),
		Location:    pgtype.Text{String: event.Location, Valid: event.Location != ""},
		Deleted:     event.Deleted,
	}); err != nil {
		return fmt.Errorf("update event %+v | %w", e, err)
	}

	return nil
}

func (e *Event) Delete(ctx context.Context, eventID int) error {
	if err := e.repo.queries(ctx).EventDelete(ctx, int32(eventID)); err != nil {
		return fmt.Errorf("delete event %d | %w", eventID, err)
	}

	return nil
}
