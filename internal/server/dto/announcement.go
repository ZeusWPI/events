package dto

import (
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
)

type Announcement struct {
	ID       int       `json:"id"`
	YearID   int       `json:"year_id" validate:"required"`
	EventIDs []int     `json:"event_ids"`
	Content  string    `json:"content" validate:"required"`
	SendTime time.Time `json:"send_time" validate:"required"`
	Send     bool      `json:"send"`
	Error    string    `json:"error,omitzero"`
}

func AnnouncementDTO(announcement *model.Announcement) Announcement {
	return Announcement(*announcement)
}

func (a *Announcement) ToModel() *model.Announcement {
	announcement := model.Announcement(*a)
	return &announcement
}
