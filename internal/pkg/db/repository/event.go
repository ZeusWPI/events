package repository

import (
	"context"
	"fmt"

	"github.com/ZeusWPI/events/internal/pkg/db/sqlc"
	"github.com/ZeusWPI/events/internal/pkg/models"
	"github.com/ZeusWPI/events/pkg/db"
	"github.com/ZeusWPI/events/pkg/util"
	"github.com/jackc/pgx/v5/pgtype"
)

// Event provides all models.Event related database operations
type Event interface {
	GetAll() ([]*models.Event, error)
	Save(*models.Event) error
}

type eventRepo struct {
	db db.DB
}

// Interface compliance
var _ Event = (*eventRepo)(nil)

// GetAll returns all events
func (r *eventRepo) GetAll() ([]*models.Event, error) {
	eventsDB, err := r.db.Queries().GetAllEvents(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Unable to get all events | %w", err)
	}

	return util.SliceMap(
			eventsDB,
			func(ev sqlc.Event) *models.Event { return sqlcEventToModel(ev) }),
		nil
}

// Save saves a new event
func (r *eventRepo) Save(e *models.Event) error {
	var eventDB sqlc.Event
	var err error

	if e.ID == 0 {
		// Create
		eventDB, err = r.db.Queries().CreateEvent(context.Background(), sqlc.CreateEventParams{
			Url:          e.URL,
			Name:         e.Name,
			Description:  pgtype.Text{String: e.Description, Valid: true},
			StartTime:    pgtype.Timestamptz{Time: e.StartTime, Valid: true},
			EndTime:      pgtype.Timestamptz{Time: e.EndTime, Valid: true},
			AcademicYear: e.AcademicYear,
			Location:     pgtype.Text{String: e.Location, Valid: true},
		})
	} else {
		// Update
		eventDB, err = r.db.Queries().UpdateEvent(context.Background(), sqlc.UpdateEventParams{
			Url:          e.URL,
			Name:         e.Name,
			Description:  pgtype.Text{String: e.Description, Valid: true},
			StartTime:    pgtype.Timestamptz{Time: e.StartTime, Valid: true},
			EndTime:      pgtype.Timestamptz{Time: e.EndTime, Valid: true},
			AcademicYear: e.AcademicYear,
			Location:     pgtype.Text{String: e.Location, Valid: true},
		})
	}

	if err != nil {
		return fmt.Errorf("Unable to save event %+v | %w", *e, err)
	}

	*e = *sqlcEventToModel(eventDB)

	return nil
}

func sqlcEventToModel(e sqlc.Event) *models.Event {
	return &models.Event{
		ID:           int(e.ID),
		URL:          e.Url,
		Name:         e.Name,
		Description:  e.Description.String,
		StartTime:    e.StartTime.Time,
		EndTime:      e.EndTime.Time,
		AcademicYear: e.AcademicYear,
		Location:     e.Location.String,
		CreatedAt:    e.CreatedAt.Time,
		UpdatedAt:    e.UpdatedAt.Time,
	}
}
