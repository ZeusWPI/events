package mattermost

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ZeusWPI/events/internal/check"
	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/task"
	"github.com/ZeusWPI/events/pkg/utils"
)

const announcementTask = "Announcement send"

func (m *Mattermost) sendAnnouncement(ctx context.Context, announcement model.Announcement) error {
	if err := m.SendMessage(ctx, Message{
		ChannelID: m.announcementChannel,
		Message:   announcement.Content,
	}); err != nil {
		announcement.Error = err.Error()
		if dbErr := m.Announcement.repoAnnouncement.Error(ctx, announcement); dbErr != nil {
			err = errors.Join(err, dbErr)
		}

		return err
	}

	if err := m.Announcement.repoAnnouncement.Send(ctx, announcement.ID); err != nil {
		return err
	}

	return nil
}

// If an announcement is already scheduled then update needs to be set to true so that it cancels it first
func (m *Mattermost) ScheduleAnnouncement(ctx context.Context, announcement model.Announcement, update bool) error {
	name := fmt.Sprintf("%s %d", announcementTask, announcement.ID)

	if update {
		// It's fine if it errors
		_ = m.task.RemoveOnceByName(name)
	}

	if announcement.SendTime.Before(time.Now()) {
		announcement.Error = "Announcement send time was in the past"
		if err := m.Announcement.repoAnnouncement.Error(ctx, announcement); err != nil {
			return err
		}

		return nil
	}

	interval := time.Until(announcement.SendTime)

	if err := m.task.AddOnce(task.NewTask(
		name,
		interval,
		func(ctx context.Context) error { return m.sendAnnouncement(ctx, announcement) },
	)); err != nil {
		return err
	}

	return nil
}

// Struct to adhere to the check interface

type announcement struct {
	repoAnnouncement repository.Announcement
}

func newAnnouncement(repo repository.Repository) *announcement {
	return &announcement{
		repoAnnouncement: *repo.NewAnnouncement(),
	}
}

// Interface compliance
var _ check.Check = (*announcement)(nil)

func (a *announcement) Description() string {
	return "Write a Mattermost announcement"
}

func (a *announcement) Status(ctx context.Context, events []model.Event) []check.StatusResult {
	statusses := make(map[int]check.StatusResult)
	for _, event := range events {
		statusses[event.ID] = check.StatusResult{
			EventID: event.ID,
			Done:    false,
			Error:   nil,
		}
	}

	announcements, err := a.repoAnnouncement.GetByEvents(ctx, events)
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
