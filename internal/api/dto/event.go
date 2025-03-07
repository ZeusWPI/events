package dto

import (
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/pkg/util"
)

// Event is the data transferable object version of the model Event
type Event struct {
	ID          int         `json:"id"`
	URL         string      `json:"url"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	StartTime   time.Time   `json:"start_time"`
	EndTime     *time.Time  `json:"end_time,omitempty"` // Pointer to support omitempty
	Location    string      `json:"location"`
	Year        Year        `json:"year"`
	Organizers  []Organizer `json:"organizers"`
}

// EventDTO converts a model Event to a DTO Event
func EventDTO(e *model.Event) Event {
	endTime := &e.EndTime
	if e.EndTime.IsZero() {
		endTime = nil
	}

	organizers := make([]Organizer, 0, len(e.Organizers))
	for _, o := range e.Organizers {
		organizers = append(organizers, Organizer{
			ID:   o.ID,
			Role: o.Role,
			Name: o.Member.Name,
		})
	}

	return Event{
		ID:          e.ID,
		URL:         e.URL,
		Name:        e.Name,
		Description: e.Description,
		StartTime:   e.StartTime,
		EndTime:     endTime,
		Location:    e.Location,
		Year: Year{
			ID:        e.Year.ID,
			StartYear: e.Year.StartYear,
			EndYear:   e.Year.EndYear,
		},
		Organizers: organizers,
	}
}

// ToModel coverts a dto Event to a model Event
func (e *Event) ToModel() *model.Event {
	endTime := time.Time{}
	if e.EndTime != nil {
		endTime = *e.EndTime
	}

	return &model.Event{
		ID:          e.ID,
		URL:         e.URL,
		Name:        e.Name,
		Description: e.Description,
		StartTime:   e.StartTime,
		EndTime:     endTime,
		Location:    e.Location,
		Year:        *e.Year.ToModel(),
		Organizers: util.SliceMap(e.Organizers, func(o Organizer) model.Board {
			return model.Board{
				ID:   o.ID,
				Role: o.Role,
				Member: model.Member{
					Name: o.Name,
				},
				Year: *e.Year.ToModel(),
			}
		}),
	}
}
