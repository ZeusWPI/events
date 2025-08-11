package model

import (
	"time"

	"github.com/ZeusWPI/events/internal/db/sqlc"
	"github.com/ZeusWPI/events/pkg/utils"
)

type Announcement struct {
	ID       int
	YearID   int
	EventIDs []int
	AuthorID int
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
		AuthorID: int(announcement.AuthorID),
		Content:  announcement.Content,
		SendTime: announcement.SendTime.Time,
		Send:     announcement.Send,
		Error:    err,
	}
}

func AnnouncementEventsModel(announcements []sqlc.AnnouncementGetByIDRow) []*Announcement {
	announcementMap := make(map[int32]*Announcement)

	for _, a := range announcements {
		if _, ok := announcementMap[a.ID]; !ok {
			err := ""
			if a.Error.Valid {
				err = a.Error.String
			}
			announcementMap[a.ID] = &Announcement{
				ID:       int(a.ID),
				YearID:   int(a.YearID),
				EventIDs: []int{},
				AuthorID: int(a.AuthorID),
				Content:  a.Content,
				SendTime: a.SendTime.Time,
				Send:     a.Send,
				Error:    err,
			}
		}

		if a.EventID.Valid {
			announcementMap[a.ID].EventIDs = append(announcementMap[a.ID].EventIDs, int(a.EventID.Int32))
		}
	}

	return utils.MapValues(announcementMap)
}
