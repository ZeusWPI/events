package dto

import (
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
)

type Check struct {
	ID          int               `json:"id"`
	EventID     int               `json:"event_id" validate:"required"`
	Status      model.CheckStatus `json:"status"`
	Message     string            `json:"message,omitzero"`
	Description string            `json:"description" validate:"required"`
	Deadline    time.Duration     `json:"deadline,omitempty"`
	Type        model.CheckType   `json:"type"`
	CreatorID   int               `json:"creator_id,omitzero"`
}

func CheckDTO(check *model.Check) Check {
	return Check{
		ID:          check.ID,
		EventID:     check.EventID,
		Status:      check.Status,
		Message:     check.Message,
		Description: check.Description,
		Deadline:    check.Deadline,
		Type:        check.Type,
		CreatorID:   check.CreatorID,
	}
}

type CheckUpdate struct {
	ID          int               `json:"id" validate:"required"`
	Status      model.CheckStatus `json:"status" validate:"required"`
	Description string            `json:"description" validate:"required"`
}
