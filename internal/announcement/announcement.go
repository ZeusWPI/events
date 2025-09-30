// Package announcement sends announcements to mattermost
package announcement

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/ZeusWPI/events/internal/check"
	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/task"
	"github.com/ZeusWPI/events/pkg/config"
	"github.com/ZeusWPI/events/pkg/mattermost"
)

const (
	checkUID = "check-mattermost"
	taskUID  = "task-mattermost"
)

type Client struct {
	development         bool
	announcementChannel string
	deadline            time.Duration

	repoAnnouncement repository.Announcement
	repoEvent        repository.Event
	m                mattermost.Client
}

func New(repo repository.Repository) (*Client, error) {
	announcementChannel := config.GetDefaultString("announcement.channel", "")
	if announcementChannel == "" {
		return nil, errors.New("no mattermost announcement channel id set")
	}

	m, err := mattermost.New()
	if err != nil {
		return nil, err
	}

	client := &Client{
		development:         config.IsDev(),
		announcementChannel: announcementChannel,
		deadline:            config.GetDefaultDuration("announcement.deadline_s", 3*24*60*60),
		repoAnnouncement:    *repo.NewAnnouncement(),
		repoEvent:           *repo.NewEvent(),
		m:                   *m,
	}

	// Register check
	if err := check.Manager.Register(context.Background(), check.NewCheck(
		checkUID,
		"Write a Mattermost announcement",
		client.deadline,
	)); err != nil {
		return nil, err
	}

	// Reschedule announcements
	if err := client.startup(context.Background()); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) startup(ctx context.Context) error {
	// Reschedule all announcements
	announcements, err := c.repoAnnouncement.GetUnsend(ctx)
	if err != nil {
		return err
	}

	for _, announcement := range announcements {
		if err := c.ScheduleAnnouncement(ctx, *announcement); err != nil {
			return err
		}
	}

	return nil
}

// Create handles a new announcement
// It changes check statusses and schedules the announcement
func (c *Client) Create(ctx context.Context, announcement model.Announcement) error {
	// Schedule announcement
	if err := c.ScheduleAnnouncement(ctx, announcement); err != nil {
		return fmt.Errorf("schedule announcement %+v | %w", announcement, err)
	}

	// Update checks
	if err := c.handleEvent(ctx, announcement.EventIDs, true); err != nil {
		return fmt.Errorf("process announcement create %+v | %w", announcement, err)
	}

	return nil
}

// Update handles an update to an announcement
// It will change the check statusses and reschedule the announcement
func (c *Client) Update(ctx context.Context, oldAnnouncement, newAnnouncement model.Announcement) error {
	// Remove old scheduled announcement
	if err := task.Manager.RemoveByUID(getTaskUID(oldAnnouncement)); err != nil {
		return fmt.Errorf("remove old announcement task %+v | %w", oldAnnouncement, err)
	}

	// Schedule new announcement
	if err := c.ScheduleAnnouncement(ctx, newAnnouncement); err != nil {
		return fmt.Errorf("schedule announcement %+v | %w", newAnnouncement, err)
	}

	// Update checks
	var newEvents []int
	var oldEvents []int

	for _, id := range newAnnouncement.EventIDs {
		if idx := slices.Index(oldAnnouncement.EventIDs, id); idx == -1 {
			// Event id is not present in the old announcement
			newEvents = append(newEvents, id)
		}
	}

	for _, id := range oldAnnouncement.EventIDs {
		if idx := slices.Index(newAnnouncement.EventIDs, id); idx == -1 {
			// Event id is removed in the new announcement
			oldEvents = append(oldEvents, id)
		}
	}

	if err := c.handleEvent(ctx, newEvents, true); err != nil {
		return fmt.Errorf("process announcement update added part %+v | %+v | %w", oldAnnouncement, newAnnouncement, err)
	}

	if err := c.handleEvent(ctx, oldEvents, false); err != nil {
		return fmt.Errorf("process announcement update deleted part %+v | %+v | %w", oldAnnouncement, newAnnouncement, err)
	}

	return nil
}

// Delete handles an announcement delete
// It will change the check statusses and cancel the announcement task
func (c *Client) Delete(ctx context.Context, announcement model.Announcement) error {
	// Remove scheduled task
	if err := task.Manager.RemoveByUID(getTaskUID(announcement)); err != nil {
		return fmt.Errorf("remove deleted announcement task %+v | %w", announcement, err)
	}

	// Remove the checks
	if err := c.handleEvent(ctx, announcement.EventIDs, false); err != nil {
		return fmt.Errorf("process announcement update %+v | %w", announcement, err)
	}

	return nil
}

func (c *Client) handleEvent(ctx context.Context, eventIDs []int, added bool) error {
	for _, eventID := range eventIDs {
		event, err := c.repoEvent.GetByID(ctx, eventID)
		if err != nil {
			return err
		}
		if event == nil {
			return fmt.Errorf("unknown event with id %d", eventID)
		}

		var status model.CheckStatus
		if time.Now().Add(c.deadline).Before(event.StartTime) {
			// Check is in time
			status = model.CheckDone
			if !added {
				status = model.CheckTODO
			}
		} else {
			status = model.CheckDoneLate
			if !added {
				status = model.CheckTODOLate
			}
		}

		if err := check.Manager.Update(ctx, checkUID, check.Update{
			Status:  status,
			EventID: event.ID,
		}); err != nil {
			return err
		}
	}

	return nil
}
