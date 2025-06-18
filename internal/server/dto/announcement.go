package dto

import (
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
)

type Announcement struct {
	ID       int       `json:"id"`
	EventID  int       `json:"event_id" validate:"required"`
	Content  string    `json:"content" validate:"required"`
	SendTime time.Time `json:"send_time" validate:"required"`
	Send     bool      `json:"send"`
}

func AnnouncementDTO(announcement model.Announcement) Announcement {
	return Announcement(announcement)
}

func (a *Announcement) ToModel() model.Announcement {
	return model.Announcement(*a)
}
