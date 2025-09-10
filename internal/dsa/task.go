package dsa

import (
	"context"
	"fmt"

	"github.com/ZeusWPI/events/internal/check"
	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/pkg/utils"
)

func (c *Client) Sync(ctx context.Context) error {
	activities, err := c.getActivities(ctx)
	if err != nil {
		return fmt.Errorf("get all activities %w", err)
	}

	events, err := c.repoEvent.GetFuture(ctx)
	if err != nil {
		return err
	}

	dsas, err := c.repoDSA.GetByEvents(ctx, utils.SliceDereference(events))
	if err != nil {
		return err
	}

	// In a perfect world everything is always synced
	// But to be sure we're going to check if that is still the case

	var uncreated []model.Event // Events that we haven't created on the DSA panel
	var updated []model.Event   // Events that have different info on the panel then in our db
	var deleted []model.Event   // Events that we have created but were manually deleted

	for _, event := range events {
		// If we've created it on the DSA panel then we have a DSA db entry
		dsa, ok := utils.SliceFind(dsas, func(d *model.DSA) bool { return d.EventID == event.ID })
		if !ok {
			// We haven't created it yet
			// That's weird, the create function must have failed
			uncreated = append(uncreated, *event)
			continue
		}

		// Now we know that we have created it in the past
		// But is it still on the DSA panel?
		activity, ok := utils.SliceFind(activities, func(a activity) bool { return a.ID == dsa.DsaID })
		if !ok {
			// The event got manually deleted!
			// Hmmmmmmm :thinking_face:
			if dsa.Deleted {
				// We already knew it was deleted
				continue
			}
			deleted = append(deleted, *event)
		}

		// The event is on the DSA panel, all is well in the world
		// Let's just check if the information is correct
		if !event.StartTime.Equal(activity.StartTime) || !event.EndTime.Equal(activity.EndTime) || event.Location != activity.Location || event.Description != activity.Description {
			// Someone manually updated it on the DSA panel
			// Don't overwrite it but let the user in the frontend know
			updated = append(updated, *event)
		}
	}

	// Let's now go over category and fix the issue

	for _, e := range uncreated {
		if err := c.Create(ctx, e); err != nil {
			return err
		}
	}

	for _, e := range updated {
		if err := check.Manager.Update(ctx, checkUID, check.Update{
			Status:  model.CheckWarning,
			Message: "DSA entry differs!",
			EventID: e.ID,
		}); err != nil {
			return err
		}
	}

	for _, e := range deleted {
		// Update the status
		if err := check.Manager.Update(ctx, checkUID, check.Update{
			Status:  model.CheckWarning,
			Message: "DSA entry got manually deleted",
			EventID: e.ID,
		}); err != nil {
			return err
		}

		// Update the DB entry so that we know it got deleted
		dsa, err := c.repoDSA.GetByEvent(ctx, e.ID)
		if err != nil {
			return err
		}

		dsa.Deleted = true
		if err := c.repoDSA.Update(ctx, *dsa); err != nil {
			return fmt.Errorf("update dsa entry to deleted %+v | %w", *dsa, err)
		}
	}

	return nil
}
