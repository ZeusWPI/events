package model

import "github.com/ZeusWPI/events/internal/db/sqlc"

type MailEvent struct {
	ID      int
	MailID  int
	EventID int
}

func MailEventModel(mailEvent sqlc.MailEvent) *MailEvent {
	return &MailEvent{
		ID:      int(mailEvent.ID),
		MailID:  int(mailEvent.MailID),
		EventID: int(mailEvent.EventID),
	}
}
