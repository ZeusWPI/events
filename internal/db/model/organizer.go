package model

import (
	"github.com/ZeusWPI/events/internal/db/sqlc"
)

type Organizer struct {
	ID      int
	EventID int
	BoardID int
	// Non DB fields
	Event Event
	Board Board
}

func OrganizerModel(organizer sqlc.Organizer, board sqlc.Board, event sqlc.Event, member sqlc.Member, year sqlc.Year) *Organizer {
	return &Organizer{
		ID:      int(organizer.ID),
		EventID: int(organizer.EventID),
		BoardID: int(organizer.BoardID),
		Event:   *EventModel(event, year),
		Board:   *BoardModel(board, member, year),
	}
}
