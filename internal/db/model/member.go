package model

import (
	"github.com/ZeusWPI/events/internal/db/sqlc"
)

type Member struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	ZauthID  int    `json:"zauth_id"`
}

func MemberModel(member sqlc.Member) *Member {
	username := member.Username.String
	if !member.Username.Valid {
		username = ""
	}
	zauthID := int(member.ZauthID.Int32)
	if !member.ZauthID.Valid {
		zauthID = 0
	}

	return &Member{
		ID:       int(member.ID),
		Name:     member.Name,
		Username: username,
		ZauthID:  zauthID,
	}
}

func (m *Member) Equal(m2 Member) bool {
	return m.Name == m2.Name
}
