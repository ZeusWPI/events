package mattermost

import (
	"context"
	"errors"

	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/task"
	"github.com/ZeusWPI/events/pkg/config"
)

type Mattermost struct {
	token               string
	url                 string
	announcementChannel string

	repoAnnouncement repository.Announcement
	task             *task.Manager
}

func New(repo repository.Repository, task *task.Manager) (*Mattermost, error) {
	token := config.GetDefaultString("mattermost.token", "")
	if token == "" {
		return nil, errors.New("no mattermost token set")
	}

	url := config.GetDefaultString("mattermost.url", "")
	if url == "" {
		return nil, errors.New("no mattermost url set")
	}

	announcementChannel := config.GetDefaultString("mattermost.announcement_channel", "")
	if announcementChannel == "" {
		return nil, errors.New("no mattermost announcement channel id set")
	}

	mattermost := &Mattermost{
		token:               token,
		url:                 url,
		announcementChannel: announcementChannel,
		repoAnnouncement:    *repo.NewAnnouncement(),
		task:                task,
	}

	if err := mattermost.startup(context.Background()); err != nil {
		return nil, err
	}

	return mattermost, nil
}

func (m *Mattermost) startup(ctx context.Context) error {
	// Reschedule all announcements
	announcements, err := m.repoAnnouncement.GetUnsend(ctx)
	if err != nil {
		return err
	}

	for _, announcement := range announcements {
		if err := m.ScheduleAnnouncement(ctx, *announcement, false); err != nil {
			return err
		}
	}

	return nil
}
