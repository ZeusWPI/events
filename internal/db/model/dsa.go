package model

import "github.com/ZeusWPI/events/internal/db/sqlc"

type DSA struct {
	ID      int
	EventID int
	Entry   bool
}

func DSAModel(dsa sqlc.Dsa) *DSA {
	return &DSA{
		ID:      int(dsa.ID),
		EventID: int(dsa.EventID),
		Entry:   dsa.Entry,
	}
}
