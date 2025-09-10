package dto

import (
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
)

type Mail struct {
	ID       int       `json:"id"`
	YearID   int       `json:"year_id" validate:"required"`
	EventIDs []int     `json:"event_ids"`
	AuthorID int       `json:"author_id"`
	Title    string    `json:"title" validate:"required"`
	Content  string    `json:"content" validate:"required"`
	SendTime time.Time `json:"send_time" validate:"required"`
	Send     bool      `json:"send"`
	Error    string    `json:"error,omitempty"`
}

func MailDTO(mail *model.Mail) Mail {
	return Mail(*mail)
}

func (m *Mail) ToModel() *model.Mail {
	mail := model.Mail(*m)
	return &mail
}
