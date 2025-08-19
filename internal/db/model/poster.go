package model

import (
	"github.com/ZeusWPI/events/internal/db/sqlc"
)

type Poster struct {
	ID      int    `json:"id"`
	EventID int    `json:"event_id"`
	FileID  string `json:"file_id"`
	WebpID  string `json:"webp_id"`
	SCC     bool   `json:"scc"`
}

func PosterModel(poster sqlc.Poster) *Poster {
	return &Poster{
		ID:      int(poster.ID),
		EventID: int(poster.EventID),
		FileID:  poster.FileID,
		WebpID:  poster.WebpID,
		SCC:     poster.Scc,
	}
}

func (p *Poster) Equal(p2 Poster) bool {
	return p.EventID == p2.EventID && p.FileID == p2.FileID && p.SCC == p2.SCC
}

// EqualEntry returns true if both poster are for the same event and of the same type (scc)
func (p *Poster) EqualEntry(p2 Poster) bool {
	return p.EventID == p2.EventID && p.SCC == p2.SCC
}
