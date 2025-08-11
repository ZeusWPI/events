package dto

import "github.com/ZeusWPI/events/internal/db/model"

type Organizer struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Role    string `json:"role"`
	ZauthID int    `json:"zauth_id,omitzero"`
}

func OrganizerDTO(b *model.Board) Organizer {
	return Organizer{
		ID:      b.ID,
		Name:    b.Member.Name,
		Role:    b.Role,
		ZauthID: b.Member.ZauthID,
	}
}
