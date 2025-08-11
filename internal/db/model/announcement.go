package model

import (
	"time"

	"github.com/ZeusWPI/events/internal/db/sqlc"
)

type Announcement struct {
	ID       int
	YearID   int
	EventIDs []int
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
		YearID:   int(announcement.YearID),
		EventIDs: []int{},
		Content:  announcement.Content,
		SendTime: announcement.SendTime.Time,
		Send:     announcement.Send,
		Error:    err,
	}
}
