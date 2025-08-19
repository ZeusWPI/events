package dto

import "github.com/ZeusWPI/events/internal/db/model"

type Poster struct {
	ID      int    `json:"id"`
	EventID int    `json:"event_id"`
	FileID  string `json:"-"`
	WebpID  string `json:"-"`
	SCC     bool   `json:"scc"`
}

func PosterDTO(poster *model.Poster) Poster {
	return Poster(*poster)
}

type PosterSave struct {
	ID      int    `form:"id"`
	EventID int    `form:"event_id" validate:"required"`
	SCC     bool   `form:"scc"`
	File    []byte `validate:"required,min=1"`
}
