package model

import (
	"time"

	"github.com/ZeusWPI/events/internal/db/sqlc"
)

type Announcement struct {
	ID       int
	EventID  int
	Content  string
	SendTime time.Time
	Send     bool
	Error    string
}

func AnnouncementModel(announcement sqlc.Announcement) *Announcement {
	err := ""
	if announcement.Error.Valid {
		err = announcement.Error.String
	}

	return &Announcement{
		ID:       int(announcement.ID),
		EventID:  int(announcement.EventID),
		Content:  announcement.Content,
		SendTime: announcement.SendTime.Time,
		Send:     announcement.Send,
		Error:    err,
	}
}
