package mattermost

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/task"
)

const announcementTask = "Announcement send"

func (c *Client) sendAnnouncement(ctx context.Context, announcement model.Announcement) error {
	if err := c.SendMessage(ctx, Message{
		ChannelID: c.announcementChannel,
		Message:   announcement.Content,
	}); err != nil {
		announcement.Error = err.Error()
		if dbErr := c.repoAnnouncement.Error(ctx, announcement); dbErr != nil {
			err = errors.Join(err, dbErr)
		}

		return err
	}

	if err := c.repoAnnouncement.Send(ctx, announcement.ID); err != nil {
		return err
	}

	return nil
}

// ScheduleAnnouncement schedules a new announcement to be sent
// If an announcement is already scheduled then update needs to be set to true so that it cancels it first
func (c *Client) ScheduleAnnouncement(ctx context.Context, announcement model.Announcement, update bool) error {
	name := fmt.Sprintf("%s %d", announcementTask, announcement.ID)

	if update {
		// It's fine if it errors
		_ = c.task.RemoveOnceByName(name)
	}

	if announcement.SendTime.Before(time.Now()) {
		announcement.Error = "Announcement send time was in the past"
		if err := c.repoAnnouncement.Error(ctx, announcement); err != nil {
			return err
		}

		return nil
	}

	interval := time.Until(announcement.SendTime)

	if err := c.task.AddOnce(task.NewTask(
		name,
		interval,
		func(ctx context.Context) error { return c.sendAnnouncement(ctx, announcement) },
	)); err != nil {
		return err
	}

	return nil
}
