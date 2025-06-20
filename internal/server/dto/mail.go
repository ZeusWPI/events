package dto

import (
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
)

type Mail struct {
	ID       int       `json:"id"`
	Content  string    `json:"content"`
	SendTime time.Time `json:"send_time"`
	Send     bool      `json:"send"`
	Error    string    `json:"error,omitzero"`
}

func MailDTO(mail *model.Mail) Mail {
	return Mail(*mail)
}

func (m *Mail) ToModel() model.Mail {
	return model.Mail(*m)
}

type MailSave struct {
	ID       int       `json:"id"`
	Content  string    `json:"content" validate:"required"`
	SendTime time.Time `json:"send_time" validate:"required"`
	EventIDs []int     `json:"event_ids" validate:"required"`
}
