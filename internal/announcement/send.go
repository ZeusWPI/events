package announcement

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/task"
	"github.com/ZeusWPI/events/pkg/mattermost"
)

func getTaskUID(announcement model.Announcement) string {
	return fmt.Sprintf("%s-%d", taskUID, announcement.ID)
}

func getTaskName(announcement model.Announcement) string {
	return fmt.Sprintf("Send announcement %d", announcement.ID)
}

func (c *Client) sendAnnouncement(ctx context.Context, announcement model.Announcement) error {
	if err := c.m.SendMessage(ctx, mattermost.Message{
		ChannelID: c.announcementChannel,
		Message:   announcement.Content,
	}); err != nil {
		announcement.Error = err.Error()
		if dbErr := c.repoAnnouncement.Update(ctx, announcement); dbErr != nil {
			err = errors.Join(err, dbErr)
		}

		return err
	}

	if err := c.repoAnnouncement.Send(ctx, announcement.ID); err != nil {
		return err
	}

	return nil
}

func (c *Client) ScheduleAnnouncement(ctx context.Context, announcement model.Announcement) error {
	name := getTaskName(announcement)

	if announcement.SendTime.Before(time.Now()) {
		announcement.Error = "Announcement send time was in the past"
		if err := c.repoAnnouncement.Update(ctx, announcement); err != nil {
			return err
		}

		return nil
	}

	interval := time.Until(announcement.SendTime)

	if err := task.Manager.AddOnce(task.NewTask(
		getTaskUID(announcement),
		name,
		interval,
		func(ctx context.Context) error { return c.sendAnnouncement(ctx, announcement) },
	)); err != nil {
		return err
	}

	return nil
}
