package dto

import (
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
)

type CheckType string

const (
	Automatic CheckType = "automatic"
	Manual    CheckType = "manual"
)

type Check struct {
	ID          int                   `json:"id"`
	Description string                `json:"description"`
	Status      model.CheckStatusEnum `json:"status"`
	Type        CheckType             `json:"type"`

	// Automatic check fields
	Deadline time.Duration `json:"duration,omitzero"`
	Message  string        `json:"message,omitzero"`

	// Manual check fields
	Creator Organizer `json:"creator,omitzero"`
}

type CheckCreate struct {
	EventID     int    `json:"event_id" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func (c *CheckCreate) ToModel() *model.CheckCustom {
	return &model.CheckCustom{
		EventID:     c.EventID,
		Description: c.Description,
	}
}

type CheckUpdate struct {
	ID          int                   `json:"id" validate:"required"`
	EventID     int                   `json:"event_id"`
	Description string                `json:"description"`
	Status      model.CheckStatusEnum `json:"status"`
}

func (c *CheckUpdate) ToModel() *model.CheckCustom {
	return &model.CheckCustom{
		ID:          c.ID,
		EventID:     c.EventID,
		Description: c.Description,
		Status:      c.Status,
	}
}
