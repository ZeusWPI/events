package website

import (
	"context"
	"fmt"
	"slices"

	"github.com/ZeusWPI/events/internal/check"
	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/pkg/utils"
)

const eventURL = "https://api.github.com/repos/ZeusWPI/zeus.ugent.be/contents/content/events"

func (c *Client) SyncEvents(ctx context.Context) error {
	websiteEvents, err := c.getEvents(ctx)
	if err != nil {
		return err
	}

	dbYears, err := c.yearRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	dbEvents, err := c.eventRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	// Create / update events
	for _, event := range websiteEvents {
		if exists := slices.ContainsFunc(dbEvents, func(e *model.Event) bool { return e.Equal(event) }); exists {
			// Exact copy already exists
			continue
		}

		// Is it an update?
		if oldEvent, ok := utils.SliceFind(dbEvents, func(e *model.Event) bool { return e.EqualEntry(event) }); ok {
			// Both the website and the local database contain this event
			// but they differ slightly, let's bring it up to date.
			// This situation can happen if
			//   - The website entry changed (e.g. new description)
			event.ID = oldEvent.ID
			event.YearID = oldEvent.YearID
			if err := c.dsa.Update(ctx, event); err != nil {
				return fmt.Errorf("updating dsa entry for old event %+v | %w", *oldEvent, err)
			}

			if err := c.eventRepo.Update(ctx, event); err != nil {
				return fmt.Errorf("updating event entry for old event %+v | %w", *oldEvent, err)
			}

			continue
		}

		// We now know it's a new event
		// Let's create it

		// Get or create the year
		if year, ok := utils.SliceFind(dbYears, func(y *model.Year) bool { return y.Equal(event.Year) }); ok {
			event.YearID = year.ID
		} else {
			if err := c.yearRepo.Create(ctx, &event.Year); err != nil {
				return fmt.Errorf("creating new year for new event %+v | %w", event, err)
			}
			dbYears = append(dbYears, &event.Year) // Update our db list to avoid creating duplicate years
			event.YearID = event.Year.ID
		}

		if err := c.eventRepo.Create(ctx, &event); err != nil {
			return err
		}

		if err := check.Manager.NewEvent(ctx, event); err != nil {
			return err
		}

		if err := c.dsa.Create(ctx, event); err != nil {
			return err
		}
	}

	// Refresh events
	dbEvents, err = c.eventRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	// Delete old events
	for _, event := range dbEvents {
		if ok := slices.ContainsFunc(websiteEvents, func(e model.Event) bool { return e.Equal(*event) }); !ok {
			if err := c.dsa.Delete(ctx, *event); err != nil {
				return err
			}

			if err := c.eventRepo.Delete(ctx, event.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
