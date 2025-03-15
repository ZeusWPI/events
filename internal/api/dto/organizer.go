package dto

import "github.com/ZeusWPI/events/internal/db/model"

// Organizer is the data transferable object of the model organizer
type Organizer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

// OrganizerDTO converts a model Board to a DTO Organizer
func OrganizerDTO(b *model.Board) Organizer {
	return Organizer{
		ID:   b.Member.ID,
		Name: b.Member.Name,
		Role: b.Role,
	}
}
