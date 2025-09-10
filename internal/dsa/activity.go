package dsa

import (
	"context"
	"fmt"
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
)

const (
	ActivitiesTask      = "DSA activities update"
	CreateActivitesTask = "Create activites on DSA website"
)

// NewEvent processes the DSA part for a new event on the website
func (c *Client) createEvent(ctx context.Context, event model.Event) error {
	if time.Now().Add(c.deadline).After(event.StartTime) {
		return nil
	}

	create := activityCreate{
		Title:       event.Name,
		Association: c.abbreviation,
		Description: event.Description,
		EndTime:     event.EndTime,
		StartTime:   event.StartTime,
		Location:    event.Location,
		Public:      true,
		Type:        "Cultuur",
		Terrain:     "ugent",
	}

	a, err := c.createActivity(ctx, create)
	if err != nil {
		return fmt.Errorf("create activity for event on the DSA website %+v | %w", event, err)
	}

	if err := c.repoDSA.Create(ctx, &model.DSA{
		DsaID:   a.ID,
		EventID: event.ID,
	}); err != nil {
		return err
	}

	return nil
}

// UpdateEvent processes the DSA part for when an event is updated on the website
func (c *Client) updateEvent(ctx context.Context, event model.Event) error {
	if time.Now().Add(c.deadline).After(event.StartTime) {
		return nil
	}

	dsa, err := c.repoDSA.GetByEvent(ctx, event.ID)
	if err != nil {
		return err
	}
	if dsa == nil {
		// Can only happen if someone manually messed with the database
		return fmt.Errorf("who touche my spaghetti (database)\nDSA entry not found for update event %+v", event)
	}

	// We don't do PATCH requests
	public := true
	update := activityUpdate{
		Title:       event.Name,
		Association: c.abbreviation,
		Description: event.Description,
		EndTime:     event.EndTime,
		StartTime:   event.StartTime,
		Location:    event.Location,
		Public:      &public,
		Type:        "Cultuur",
		Terrain:     "ugent",
	}

	if _, err := c.updateActivity(ctx, dsa.DsaID, update); err != nil {
		return fmt.Errorf("update activity for event on the DSA website %+v | %w", event, err)
	}

	return nil
}

func (c *Client) deleteEvent(ctx context.Context, event model.Event) error {
	if time.Now().After(event.StartTime) {
		return nil
	}

	dsa, err := c.repoDSA.GetByEvent(ctx, event.ID)
	if err != nil {
		return err
	}
	if dsa == nil {
		// Can only happen if someone manually messed with the database
		return fmt.Errorf("who touche my spaghetti (database)\nDSA entry not found for delete event %+v", event)
	}

	if _, err := c.deleteActivity(ctx, dsa.DsaID); err != nil {
		return fmt.Errorf("delete activity for event on the DSA website %+v | %w", event, err)
	}

	if err := c.repoDSA.Delete(ctx, dsa.ID); err != nil {
		return err
	}

	return nil
}
