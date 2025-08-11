package model

import (
	"time"

	"github.com/ZeusWPI/events/internal/db/sqlc"
	"github.com/ZeusWPI/events/pkg/utils"
)

type Mail struct {
	ID       int
	YearID   int
	EventIDs []int
	AuthorID int
	Title    string
	Content  string
	SendTime time.Time
	Send     bool
	Error    string
}

func MailModel(mail sqlc.Mail) *Mail {
	err := ""
	if mail.Error.Valid {
		err = mail.Error.String
	}

	return &Mail{
		ID:       int(mail.ID),
		YearID:   int(mail.YearID),
		EventIDs: []int{},
		AuthorID: int(mail.AuthorID),
		Title:    mail.Title,
		Content:  mail.Content,
		SendTime: mail.SendTime.Time,
		Send:     mail.Send,
		Error:    err,
	}
}

func MailEventsModel(mails []sqlc.MailGetByIDRow) []*Mail {
	mailMap := make(map[int32]*Mail)

	for _, m := range mails {
		if _, ok := mailMap[m.ID]; !ok {
			err := ""
			if m.Error.Valid {
				err = m.Error.String
			}
			mailMap[m.ID] = &Mail{
				ID:       int(m.ID),
				YearID:   int(m.YearID),
				EventIDs: []int{},
				AuthorID: int(m.AuthorID),
				Title:    m.Title,
				Content:  m.Content,
				SendTime: m.SendTime.Time,
				Send:     m.Send,
				Error:    err,
			}
		}

		if m.EventID.Valid {
			mailMap[m.ID].EventIDs = append(mailMap[m.ID].EventIDs, int(m.EventID.Int32))
		}
	}

	return utils.MapValues(mailMap)
}
