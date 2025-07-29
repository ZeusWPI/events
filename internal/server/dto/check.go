package dto

import (
	"github.com/ZeusWPI/events/internal/check"
)

type CheckStatus string

const (
	Finished   CheckStatus = "finished"
	Unfinished CheckStatus = "unfinished"
	Warning    CheckStatus = "warning"
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
	Status      CheckStatus `json:"status"`
	Error       error       `json:"error"`
	Source      CheckSource `json:"source"`
}

func CheckDTO(check check.EventStatus) Check {
	return Check{
		ID:          check.ID,
		EventID:     check.EventID,
		Description: check.Description,
		Status:      CheckStatus(check.Status),
		Error:       check.Error,
		Source:      CheckSource(check.Source),
	}
}
