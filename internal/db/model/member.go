package model

import (
	"github.com/ZeusWPI/events/internal/db/sqlc"
)

type Member struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	Mattermost string `json:"mattermost"`
	ZauthID    int    `json:"zauth_id"`
}

func MemberModel(member sqlc.Member) *Member {
	username := ""
	if member.Username.Valid {
		username = member.Username.String
	}
	mattermost := ""
	if member.Mattermost.Valid {
		mattermost = member.Mattermost.String
	}
	zauthID := 0
	if member.ZauthID.Valid {
		zauthID = int(member.ZauthID.Int32)
	}

	return &Member{
		ID:         int(member.ID),
		Name:       member.Name,
		Username:   username,
		Mattermost: mattermost,
		ZauthID:    zauthID,
	}
}

func (m *Member) Equal(m2 Member) bool {
	return m.Name == m2.Name && m.Mattermost == m2.Mattermost
}

func (m *Member) EqualEntry(m2 Member) bool {
	return m.Name == m2.Name
}
