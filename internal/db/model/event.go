package model

import (
	"time"

	"github.com/ZeusWPI/events/internal/db/sqlc"
)

type Event struct {
	ID          int       `json:"id"`
	FileName    string    `json:"file_name"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	YearID      int       `json:"year_id"`
	Location    string    `json:"location"`
	// Non db fields
	Year       Year     `json:"year"`
	Organizers []Board  `json:"organizers"`
	Posters    []Poster `json:"posters"`
}

func EventModel(event sqlc.Event) *Event {
	description := event.Description.String
	if !event.Description.Valid {
		description = ""
	}
	endTime := event.EndTime.Time
	if !event.EndTime.Valid {
		endTime = time.Time{}
	}
	location := event.Location.String
	if !event.Location.Valid {
		location = ""
	}

	return &Event{
		ID:          int(event.ID),
		FileName:    event.FileName,
		Name:        event.Name,
		Description: description,
		StartTime:   event.StartTime.Time,
		EndTime:     endTime,
		YearID:      int(event.YearID),
		Location:    location,
	}
}

func (e *Event) Equal(e2 Event) bool {
	return e.FileName == e2.FileName &&
		e.Name == e2.Name &&
		e.Description == e2.Description &&
		e.StartTime.Equal(e2.StartTime) &&
		e.EndTime.Equal(e2.EndTime) &&
		e.Location == e2.Location &&
		e.Year.Equal(e2.Year)
}
