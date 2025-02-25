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
	EventGetAllWithYear(context.Context) ([]*model.Event, error)
	Save(context.Context, *model.Event) error
	Delete(context.Context, *model.Event) error
}

type eventRepo struct {
	repo Repository
}

// Interface compliance
var _ Event = (*eventRepo)(nil)

// GetAll returns all events
func (r *eventRepo) EventGetAllWithYear(ctx context.Context) ([]*model.Event, error) {
	events, err := r.repo.queries(ctx).EventGetAllWithYear(ctx)
	if err != nil {
		return nil, fmt.Errorf("Unable to get all events | %v", err)
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
		return fmt.Errorf("Unable to save event %+v | %v", *e, err)
	}

	e.ID = int(id)

	return nil
}

// Delete soft deletes an event
func (r *eventRepo) Delete(ctx context.Context, e *model.Event) error {
	if e.ID == 0 {
		return fmt.Errorf("Event has no ID %+v", *e)
	}

	err := r.repo.queries(ctx).EventDelete(ctx, int32(e.ID))
	if err != nil {
		return fmt.Errorf("Unable to delete event %+v | %v", *e, err)
	}

	e.DeletedAt = time.Now() // Close enough

	return nil
}
