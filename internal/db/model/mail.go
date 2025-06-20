package model

import (
	"time"

	"github.com/ZeusWPI/events/internal/db/sqlc"
)

type Mail struct {
	ID       int
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
		Content:  mail.Content,
		SendTime: mail.SendTime.Time,
		Send:     mail.Send,
		Error:    err,
	}
}
