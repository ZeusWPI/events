// Package mail stands in for sending mails to Zeus WPI users
package mail

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/ZeusWPI/events/internal/check"
	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/task"
	"github.com/ZeusWPI/events/pkg/config"
)

const (
	checkUID = "check-mail"
	taskUID  = "task-mail"
)

type Client struct {
	deadline time.Duration

	repoEvent repository.Event
	repoMail  repository.Mail
}

func New(repo repository.Repository) (*Client, error) {
	client := &Client{
		deadline:  config.GetDefaultDuration("mail.deadline_s", 0),
		repoEvent: *repo.NewEvent(),
		repoMail:  *repo.NewMail(),
	}

	// Register check
	if err := check.Manager.Register(context.Background(), check.NewCheck(
		checkUID,
		"Cover the event in an email",
		client.deadline,
	)); err != nil {
		return nil, err
	}

	// Reschedule the mails
	if err := client.startup(context.Background()); err != nil {
		return nil, err
	}

	return client, nil
}

// startup reschedules all unsend mails
func (c *Client) startup(ctx context.Context) error {
	mails, err := c.repoMail.GetUnsend(ctx)
	if err != nil {
		return err
	}

	for _, mail := range mails {
		if err := c.ScheduleMailAll(ctx, *mail); err != nil {
			return err
		}
	}

	return nil
}

// Create handles a new mail being created
// It will change the check status and schedule the mail
func (c *Client) Create(ctx context.Context, mail model.Mail) error {
	// Schedule mail
	if err := c.ScheduleMailAll(ctx, mail); err != nil {
		return fmt.Errorf("schedule mail %+v | %w", mail, err)
	}

	// Update checks
	if err := c.handleEvent(ctx, mail.EventIDs, true); err != nil {
		return fmt.Errorf("process mail create %+v | %w", mail, err)
	}

	return nil
}

// Update handles a mail being updated
// It will change the check status and reschedule the mail
func (c *Client) Update(ctx context.Context, oldMail, newMail model.Mail) error {
	// Remove old scheduled mail
	if err := task.Manager.RemoveByUID(getTaskUID(oldMail)); err != nil {
		return fmt.Errorf("remove updated mail task %+v | %w", newMail, err)
	}

	// Schedule new mail
	if err := c.ScheduleMailAll(ctx, newMail); err != nil {
		return fmt.Errorf("schedule mail task %+v | %w", newMail, err)
	}

	// Update checks
	var newEvents []int
	var oldEvents []int

	for _, id := range newMail.EventIDs {
		if idx := slices.Index(oldMail.EventIDs, id); idx == -1 {
			// Event id is not present in the old mail
			newEvents = append(newEvents, id)
		}
	}

	for _, id := range oldMail.EventIDs {
		if idx := slices.Index(newMail.EventIDs, id); idx == -1 {
			// Event id is removed in the new mail
			oldEvents = append(oldEvents, id)
		}
	}

	if err := c.handleEvent(ctx, newEvents, true); err != nil {
		return fmt.Errorf("process mail update added part %+v | %+v | %w", oldMail, newMail, err)
	}

	if err := c.handleEvent(ctx, oldEvents, false); err != nil {
		return fmt.Errorf("process mail update deleted part %+v | %+v | %w", oldMail, newMail, err)
	}

	return nil
}

// Delete handles a mail being deleted
// It will change the check status and cancel the mail task
func (c *Client) Delete(ctx context.Context, mail model.Mail) error {
	// Remove scheduled task
	if err := task.Manager.RemoveByUID(getTaskUID(mail)); err != nil {
		return fmt.Errorf("remove deleted mail task %+v | %w", mail, err)
	}

	// Remove the checks
	if err := c.handleEvent(ctx, mail.EventIDs, false); err != nil {
		return fmt.Errorf("process mail update %+v | %w", mail, err)
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
