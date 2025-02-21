package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ZeusWPI/events/internal/pkg/db/sqlc"
	"github.com/ZeusWPI/events/internal/pkg/model"
	"github.com/ZeusWPI/events/pkg/db"
	"github.com/ZeusWPI/events/pkg/util"
	"github.com/jackc/pgx/v5/pgtype"
)

// Event provides all model.Event related database operations
type Event interface {
	GetAll() ([]*model.Event, error)
	Save(*model.Event) error
	Delete(*model.Event) error
}

type eventRepo struct {
	db db.DB

	year AcademicYear
}

// Interface compliance
var _ Event = (*eventRepo)(nil)

// GetAll returns all events
func (r *eventRepo) GetAll() ([]*model.Event, error) {
	eventsDB, err := r.db.Queries().EventGetAll(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Unable to get all events | %v", err)
	}

	years, err := r.year.GetAll()
	if err != nil {
		return nil, err
	}

	events := make([]*model.Event, 0, len(eventsDB))
	for _, e := range eventsDB {
		year, ok := util.SliceFind(years, func(y *model.AcademicYear) bool { return y.ID == int(e.AcademicYear) })
		if !ok {
			continue
		}

		event := &model.Event{
			ID:           int(e.ID),
			URL:          e.Url,
			Name:         e.Name,
			Description:  e.Description.String,
			StartTime:    e.StartTime.Time,
			EndTime:      e.EndTime.Time,
			AcademicYear: *year,
			Location:     e.Location.String,
			CreatedAt:    e.CreatedAt.Time,
			UpdatedAt:    e.UpdatedAt.Time,
			DeletedAt:    e.DeletedAt.Time,
		}
		events = append(events, event)
	}

	return events, nil
}

// Save creates a new academic year or updates an existing one
func (r *eventRepo) Save(e *model.Event) error {
	var id int32
	var err error

	if e.ID == 0 {
		// Create
		id, err = r.db.Queries().EventCreate(context.Background(), sqlc.EventCreateParams{
			Url:          e.URL,
			Name:         e.Name,
			Description:  pgtype.Text{String: e.Description, Valid: true},
			StartTime:    pgtype.Timestamptz{Time: e.StartTime, Valid: true},
			EndTime:      pgtype.Timestamptz{Time: e.EndTime, Valid: true},
			AcademicYear: int32(e.AcademicYear.ID),
			Location:     pgtype.Text{String: e.Location, Valid: true},
		})
	} else {
		// Update
		id = int32(e.ID)
		err = r.db.Queries().EventUpdate(context.Background(), sqlc.EventUpdateParams{
			ID:           int32(e.ID),
			Url:          e.URL,
			Name:         e.Name,
			Description:  pgtype.Text{String: e.Description, Valid: true},
			StartTime:    pgtype.Timestamptz{Time: e.StartTime, Valid: true},
			EndTime:      pgtype.Timestamptz{Time: e.EndTime, Valid: true},
			AcademicYear: int32(e.AcademicYear.ID),
			Location:     pgtype.Text{String: e.Location, Valid: true},
		})
	}

	if err != nil {
		return fmt.Errorf("Unable to save event %+v | %v", *e, err)
	}

	e.ID = int(id)

	return nil
}

// Delete soft deletes an event
func (r *eventRepo) Delete(e *model.Event) error {
	if e.ID == 0 {
		return fmt.Errorf("Event has no ID %+v", *e)
	}

	err := r.db.Queries().EventDelete(context.Background(), int32(e.ID))
	if err != nil {
		return fmt.Errorf("Unable to delete event %+v | %v", *e, err)
	}

	e.DeletedAt = time.Now() // Close enough

	return nil
}
