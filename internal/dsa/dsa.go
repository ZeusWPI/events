// Package dsa controls the syncronization between DSA and events
package dsa

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ZeusWPI/events/internal/check"
	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/task"
	"github.com/ZeusWPI/events/pkg/config"
)

const (
	checkUID = "check-dsa"
	taskUID  = "task-dsa"
)

// Change these values if you have a local dsa instance @xander
const (
	dsaURL       = "https://dsa.ugent.be/api"
	abbreviation = "zeus"
)

type Client struct {
	development bool
	key         string
	deadline    time.Duration

	repoDSA   repository.DSA
	repoEvent repository.Event
	repoYear  repository.Year
}

func New(repo repository.Repository) (*Client, error) {
	dsaKey := config.GetDefaultString("dsa.key", "")
	if dsaKey == "" {
		return nil, errors.New("no dsa api key set")
	}

	client := &Client{
		development: config.IsDev(),
		key:         dsaKey,
		deadline:    config.GetDefaultDuration("dsa.deadline_s", 3*24*60*60),
		repoDSA:     *repo.NewDSA(),
		repoEvent:   *repo.NewEvent(),
		repoYear:    *repo.NewYear(),
	}

	// Register task
	if err := task.Manager.AddRecurring(context.Background(), task.NewTask(
		taskUID,
		"DSA events synchronization",
		config.GetDefaultDuration("dsa.syncronize_s", 24*60*60),
		client.Sync,
	)); err != nil {
		return nil, err
	}

	// Register check
	if err := check.Manager.Register(context.Background(), check.NewCheck(
		checkUID,
		"Add event to the DSA website",
		client.deadline,
	)); err != nil {
		return nil, err
	}

	return client, nil
}

// Create handles a new event
func (c *Client) Create(ctx context.Context, event model.Event) error {
	// Create on the dsa website
	if err := c.createEvent(ctx, event); err != nil {
		return fmt.Errorf("create event on the dsa website %+v | %w", event, err)
	}

	// Update checks
	if err := c.handleEvent(ctx, event, true); err != nil {
		return fmt.Errorf("update dsa checks for event %+v | %w", event, err)
	}

	return nil
}

// Update handles an update to an event
func (c *Client) Update(ctx context.Context, newEvent model.Event) error {
	// Update the dsa website
	if err := c.updateEvent(ctx, newEvent); err != nil {
		return fmt.Errorf("update event on the dsa website %+v | %w", newEvent, err)
	}

	return nil
}

// Delete handles an event delete
func (c *Client) Delete(ctx context.Context, event model.Event) error {
	// Delete on the dsa website
	if err := c.deleteEvent(ctx, event); err != nil {
		return fmt.Errorf("delete event on the dsa website %+v | %w", event, err)
	}

	// Update checks
	if err := c.handleEvent(ctx, event, false); err != nil {
		return fmt.Errorf("update dsa checks for event %+v | %w", event, err)
	}

	return nil
}

func (c *Client) handleEvent(ctx context.Context, event model.Event, added bool) error {
	var status model.CheckStatus
	var message string

	if time.Now().Add(c.deadline).Before(event.StartTime) {
		// Check is in time
		status = model.CheckDone
		if !added {
			status = model.CheckTODO
		}
	} else {
		status = model.CheckTODOLate
		message = "Too late to add to the DSA website"
	}

	if err := check.Manager.Update(ctx, checkUID, check.Update{
		Status:  status,
		Message: message,
		EventID: event.ID,
	}); err != nil {
		return err
	}

	return nil
}
