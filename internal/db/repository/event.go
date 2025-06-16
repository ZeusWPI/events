package repository

import (
	"context"
	"database/sql"
	"encoding/json"
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

func (e *Event) GetAllWithYear(ctx context.Context) ([]*model.Event, error) {
	events, err := e.repo.queries(ctx).EventGetAllWithYear(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get all events with year %w", err)
	}

	return utils.SliceMap(events, func(e sqlc.EventGetAllWithYearRow) *model.Event {
		return &model.Event{
			ID:          int(e.ID),
			FileName:    e.FileName,
			Name:        e.Name,
			Description: e.Description.String,
			StartTime:   e.StartTime.Time,
			EndTime:     e.EndTime.Time,
			YearID:      int(e.YearID),
			Location:    e.Location.String,
			Year: model.Year{
				ID:    int(e.ID_2),
				Start: int(e.YearStart),
				End:   int(e.YearEnd),
			},
			Organizers: make([]model.Board, 0),
		}
	}), nil
}

func (e *Event) GetByYearPopulated(ctx context.Context, yearID int) ([]*model.Event, error) {
	eventsDB, err := e.repo.queries(ctx).EventGetByYearPopulated(ctx, int32(yearID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get all populated events by year %d | %w", yearID, err)
	}

	events := make([]*model.Event, 0, len(eventsDB))
	for _, bytes := range eventsDB {
		var event model.Event
		if err := json.Unmarshal(bytes, &event); err != nil {
			return nil, fmt.Errorf("unmarshal event json %w", err)
		}
		events = append(events, &event)
	}

	return events, nil
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
		Description: pgtype.Text{String: event.Description, Valid: true},
		StartTime:   pgtype.Timestamptz{Time: event.StartTime, Valid: true},
		EndTime:     pgtype.Timestamptz{Time: event.EndTime, Valid: true},
		YearID:      int32(event.YearID),
		Location:    pgtype.Text{String: event.Location, Valid: true},
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
