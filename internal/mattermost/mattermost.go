package mattermost

import "github.com/ZeusWPI/events/internal/db/repository"

type Mattermost struct {
	Announcement *Announcement
}

func New(repo repository.Repository) *Mattermost {
	return &Mattermost{
		Announcement: newAnnouncement(repo),
	}
}
