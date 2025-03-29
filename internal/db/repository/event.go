package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/sqlc"
	"github.com/ZeusWPI/events/pkg/util"
	"github.com/jackc/pgx/v5/pgtype"
)

// Event provides all model.Event related database operations
type Event interface {
	GetAllWithYear(context.Context) ([]*model.Event, error)
	GetByYearWithAll(context.Context, model.Year) ([]*model.Event, error)
	Save(context.Context, *model.Event) error
	Delete(context.Context, *model.Event) error
}

type eventRepo struct {
	repo Repository

	organizer Organizer
}

// Interface compliance
var _ Event = (*eventRepo)(nil)

// GetAll returns all events
func (r *eventRepo) GetAllWithYear(ctx context.Context) ([]*model.Event, error) {
	events, err := r.repo.queries(ctx).EventGetAllWithYear(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get all events %w", err)
	}

	return util.SliceMap(events, func(e sqlc.EventGetAllWithYearRow) *model.Event {
		return &model.Event{
			ID:          int(e.ID),
			URL:         e.Url,
			Name:        e.Name,
			Description: e.Description.String,
			StartTime:   e.StartTime.Time,
			EndTime:     e.EndTime.Time,
			Year: model.Year{
				ID:        int(e.ID_2),
				StartYear: int(e.StartYear),
				EndYear:   int(e.EndYear),
			},
			Location:   e.Location.String,
			Organizers: make([]model.Board, 0),
			CreatedAt:  e.CreatedAt.Time,
			UpdatedAt:  e.UpdatedAt.Time,
			DeletedAt:  e.DeletedAt.Time,
		}
	}), nil
}

// GetByYearWithAll returns all events of a given year with all model.Event fields populated
func (r *eventRepo) GetByYearWithAll(ctx context.Context, year model.Year) ([]*model.Event, error) {
	eventsDB, err := r.repo.queries(ctx).EventGetByYearWithYear(ctx, int32(year.ID))
	if err != nil {
		return nil, fmt.Errorf("unable to get all events by year %+v | %w", year, err)
	}

	organizers, err := r.organizer.GetByYearWithBoard(ctx, year)
	if err != nil {
		return nil, err
	}

	events := util.SliceMap(eventsDB, func(e sqlc.EventGetByYearWithYearRow) *model.Event {
		return &model.Event{
			ID:          int(e.ID),
			URL:         e.Url,
			Name:        e.Name,
			Description: e.Description.String,
			StartTime:   e.StartTime.Time,
			EndTime:     e.EndTime.Time,
			Location:    e.Location.String,
			Year: model.Year{
				ID:        int(e.Year),
				StartYear: int(e.StartYear),
				EndYear:   int(e.EndYear),
			},
			Organizers: make([]model.Board, 0),
			CreatedAt:  e.CreatedAt.Time,
			UpdatedAt:  e.UpdatedAt.Time,
			DeletedAt:  e.DeletedAt.Time,
		}
	})

	for _, organizer := range organizers {
		for i, event := range events {
			if organizer.Event.ID == event.ID {
				events[i].Organizers = append(events[i].Organizers, organizer.Board)
				break
			}
		}
	}

	return events, nil
}

// Save creates a new event or updates an existing one
func (r *eventRepo) Save(ctx context.Context, e *model.Event) error {
	var id int32
	var err error

	if e.ID == 0 {
		// Create
		id, err = r.repo.queries(ctx).EventCreate(ctx, sqlc.EventCreateParams{
			Url:         e.URL,
			Name:        e.Name,
			Description: pgtype.Text{String: e.Description, Valid: true},
			StartTime:   pgtype.Timestamptz{Time: e.StartTime, Valid: true},
			EndTime:     pgtype.Timestamptz{Time: e.EndTime, Valid: true},
			Year:        int32(e.Year.ID),
			Location:    pgtype.Text{String: e.Location, Valid: true},
		})
	} else {
		// Update
		id = int32(e.ID)
		err = r.repo.queries(ctx).EventUpdate(ctx, sqlc.EventUpdateParams{
			ID:          int32(e.ID),
			Url:         e.URL,
			Name:        e.Name,
			Description: pgtype.Text{String: e.Description, Valid: true},
			StartTime:   pgtype.Timestamptz{Time: e.StartTime, Valid: true},
			EndTime:     pgtype.Timestamptz{Time: e.EndTime, Valid: true},
			Year:        int32(e.Year.ID),
			Location:    pgtype.Text{String: e.Location, Valid: true},
		})
	}

	if err != nil {
		return fmt.Errorf("unable to save event %+v | %w", *e, err)
	}

	e.ID = int(id)

	return nil
}

// Delete soft deletes an event
func (r *eventRepo) Delete(ctx context.Context, e *model.Event) error {
	if e.ID == 0 {
		return fmt.Errorf("Event has no ID %+v", *e)
	}

	if err := r.repo.queries(ctx).EventDelete(ctx, int32(e.ID)); err != nil {
		return fmt.Errorf("unable to delete event %+v | %w", *e, err)
	}

	e.DeletedAt = time.Now() // Close enough

	return nil
}
