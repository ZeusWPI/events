package dto

import (
	"fmt"
	"strings"
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/pkg/utils"
)

const eventURL = "https://zeus.gent/events"

type Event struct {
	ID           int          `json:"id"`
	URL          string       `json:"url"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	StartTime    time.Time    `json:"start_time"`
	EndTime      *time.Time   `json:"end_time,omitempty"` // Pointer to support omitempty
	Location     string       `json:"location"`
	Year         Year         `json:"year"`
	Organizers   []Organizer  `json:"organizers"`
	Checks       []Check      `json:"checks"`
	Announcement Announcement `json:"announcement,omitzero"`
}

func EventDTO(event *model.Event) Event {
	endTime := &event.EndTime
	if event.EndTime.IsZero() {
		endTime = nil
	}

	return Event{
		ID:          event.ID,
		URL:         fmt.Sprintf("%s/%s/%s", eventURL, event.Year.String(), event.FileName),
		Name:        event.Name,
		Description: event.Description,
		StartTime:   event.StartTime,
		EndTime:     endTime,
		Location:    event.Location,
		Year: Year{
			ID:    event.Year.ID,
			Start: event.Year.Start,
			End:   event.Year.End,
		},
		Organizers: utils.SliceMap(utils.SliceReference(event.Organizers), OrganizerDTO),
	}
}

func (event *Event) ToModel() *model.Event {
	endTime := time.Time{}
	if event.EndTime != nil {
		endTime = *event.EndTime
	}

	fileName := ""
	urlParts := strings.Split(event.URL, "/")
	if len(urlParts) == 3 {
		fileName = urlParts[2]
	}

	return &model.Event{
		ID:          event.ID,
		FileName:    fileName,
		Name:        event.Name,
		Description: event.Description,
		StartTime:   event.StartTime,
		EndTime:     endTime,
		Location:    event.Location,
		Year:        *event.Year.ToModel(),
		Organizers: utils.SliceMap(event.Organizers, func(o Organizer) model.Board {
			return model.Board{
				Member: model.Member{
					ID:   o.ID,
					Name: o.Name,
				},
				Year: *event.Year.ToModel(),
				Role: o.Role,
			}
		}),
	}
}

type EventOrganizers struct {
	EventID    int   `json:"event_id" validate:"required"`
	Organizers []int `json:"organizers" validate:"required"`
}
