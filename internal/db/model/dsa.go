package model

import "github.com/ZeusWPI/events/internal/db/sqlc"

type DSA struct {
	ID      int
	EventID int
	DsaID   int
}

func DSAModel(dsa sqlc.Dsa) *DSA {
	dsaID := 0

	if dsa.DsaID.Valid {
		dsaID = int(dsa.DsaID.Int32)
	}

	return &DSA{
		ID:      int(dsa.ID),
		EventID: int(dsa.EventID),
		DsaID:   dsaID,
	}
}
