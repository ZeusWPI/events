package dto

import "github.com/ZeusWPI/events/internal/db/model"

type Organizer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

func OrganizerDTO(b *model.Board) Organizer {
	return Organizer{
		ID:   b.Member.ID,
		Name: b.Member.Name,
		Role: b.Role,
	}
}
