package dto

import (
	"github.com/ZeusWPI/events/internal/check"
)

type CheckSource string

const (
	Automatic CheckSource = "automatic"
	Manual    CheckSource = "manual"
)

type Check struct {
	ID          int         `json:"id"`
	EventID     int         `json:"event_id" validate:"required"`
	Description string      `json:"description" validate:"required"`
	Done        bool        `json:"done"`
	Error       error       `json:"error"`
	Source      CheckSource `json:"source"`
}

func CheckDTO(check check.Status) Check {
	return Check{
		ID:          check.ID,
		EventID:     check.EventID,
		Description: check.Description,
		Done:        check.Done,
		Error:       check.Error,
		Source:      CheckSource(check.Source),
	}
}
