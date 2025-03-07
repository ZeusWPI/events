package dto

import "github.com/ZeusWPI/events/internal/db/model"

// Organizer is the data transferable object of the model organizer
type Organizer struct {
	ID   int    `json:"id"`
	Role string `json:"role"`
	Name string `json:"name"`
}

// OrganizerDTO converts a model Board to a DTO Organizer
func OrganizerDTO(b *model.Board) Organizer {
	return Organizer{
		ID:   b.ID,
		Role: b.Role,
		Name: b.Member.Name,
	}
}
