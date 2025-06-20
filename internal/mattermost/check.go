package mattermost

import (
	"context"

	"github.com/ZeusWPI/events/internal/check"
	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/pkg/utils"
)

type CheckAnnouncement struct {
	repoAnnouncement repository.Announcement
}

func (m *Mattermost) NewCheckAnnouncement() *CheckAnnouncement {
	return &CheckAnnouncement{
		repoAnnouncement: m.repoAnnouncement,
	}
}

// Interface compliance
var _ check.Check = (*CheckAnnouncement)(nil)

func (c *CheckAnnouncement) Description() string {
	return "Write a Mattermost announcement"
}

func (c *CheckAnnouncement) Status(ctx context.Context, events []model.Event) []check.StatusResult {
	statusses := make(map[int]check.StatusResult)
	for _, event := range events {
		statusses[event.ID] = check.StatusResult{
			EventID: event.ID,
			Done:    false,
			Error:   nil,
		}
	}

	announcements, err := c.repoAnnouncement.GetByEvents(ctx, events)
	if err != nil {
		for k, v := range statusses {
			v.Error = err
			statusses[k] = v
		}

		return utils.MapValues(statusses)
	}

	for _, announcement := range announcements {
		if status, ok := statusses[announcement.EventID]; ok {
			status.Done = true
			statusses[announcement.EventID] = status
		}
	}

	return utils.MapValues(statusses)
}
