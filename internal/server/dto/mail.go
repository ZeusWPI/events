package dto

import (
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
)

type Mail struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	SendTime time.Time `json:"send_time"`
	Send     bool      `json:"send"`
	EventIDs []int     `json:"-"`
	Error    string    `json:"error,omitzero"`
	Events   []Event   `json:"events"`
}

func MailDTO(mail *model.Mail) Mail {
	return Mail{
		ID:       mail.ID,
		Title:    mail.Title,
		Content:  mail.Content,
		SendTime: mail.SendTime,
		Send:     mail.Send,
		EventIDs: mail.EventIDs,
		Error:    mail.Error,
	}
}

func (m *Mail) ToModel() model.Mail {
	return model.Mail{
		ID:       m.ID,
		Title:    m.Title,
		Content:  m.Content,
		SendTime: m.SendTime,
		Send:     m.Send,
		EventIDs: m.EventIDs,
		Error:    m.Error,
	}
}

type MailSave struct {
	ID       int       `json:"id"`
	Title    string    `json:"title" validate:"required"`
	Content  string    `json:"content" validate:"required"`
	SendTime time.Time `json:"send_time" validate:"required"`
	EventIDs []int     `json:"event_ids" validate:"required,min=1"`
}
