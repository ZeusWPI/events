package model

import "github.com/ZeusWPI/events/internal/db/sqlc"

type Poster struct {
	ID      int    `json:"id"`
	EventID int    `json:"event_id"`
	FileID  string `json:"file_id"`
	SCC     bool   `json:"scc"`
}

func PosterModel(poster sqlc.Poster) *Poster {
	return &Poster{
		ID:      int(poster.ID),
		EventID: int(poster.EventID),
		FileID:  poster.FileID,
		SCC:     poster.Scc,
	}
}
