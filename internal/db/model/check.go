package model

import "github.com/ZeusWPI/events/internal/db/sqlc"

type Check struct {
	ID          int
	EventID     int
	Description string
	Done        bool
}

func CheckModel(check sqlc.Check) *Check {
	return &Check{
		ID:          int(check.ID),
		EventID:     int(check.EventID),
		Description: check.Description,
		Done:        check.Done,
	}
}
