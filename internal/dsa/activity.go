package dsa

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/pkg/utils"
)

const ActivitiesTask = "DSA activities update"
const CreateActivitesTask = "Create activites on DSA website"

func (d *DSA) CreateActivities(ctx context.Context) error {
	activities, err := d.getActivities(ctx)
	if err != nil {
		return fmt.Errorf("get activities: %w", err)
	}

	events, err := d.repoEvent.GetFutureWithYear(ctx)
	if err != nil {
		return fmt.Errorf("get events by year: %w", err)
	}

	closeEvents := utils.SliceFilter(events, func(e *model.Event) bool { return e.StartTime.Before(time.Now().AddDate(0, 0, 14)) })

	dsas, err := d.repoDSA.GetByEvents(ctx, utils.SliceDereference(closeEvents))
	if err != nil {
		return fmt.Errorf("get local dsa activities: %w", err)
	}

	uncreatedActivities := utils.SliceFilter(dsas, func(d *model.DSA) bool { return d.DsaID == 0 })

	for _, dsa := range uncreatedActivities {
		event, err := d.repoEvent.GetByID(ctx, dsa.EventID)
		if err != nil {
			return fmt.Errorf("get event: %w", err)
		}

		if event == nil {
			return fmt.Errorf("no event with ID %d", dsa.EventID)
		}

		act, activityOk := utils.SliceFind(activities, func(a activity) bool { return a.StartTime.Equal(event.StartTime) })

		if !activityOk {
			// Event is not yet on the dsa website

			activityCreate := activityCreate{
				Title:       event.Name,
				Association: d.abbreviation,
				Description: event.Description,
				EndTime:     event.EndTime,
				StartTime:   event.StartTime,
				Location:    event.Location,
				Public:      true,
				Type:        "Cultuur",
				Terrain:     "ugent",
			}

			response, err := d.createActivity(ctx, activityCreate)
			if err != nil {
				return fmt.Errorf("create dsa activity: %w", err)
			}

			dsa.DsaID = response.ID

		} else {
			// Event is already on the dsa website, added manually
			dsa.DsaID = act.ID
		}

		if err := d.repoDSA.Update(ctx, *dsa); err != nil {
			return fmt.Errorf("update local dsa record: %w", err)
		}
	}

	return nil
}

func (d *DSA) UpdateActivities(ctx context.Context) error {
	activities, err := d.getActivities(ctx)
	if err != nil {
		return fmt.Errorf("get dsa activities: %w", err)
	}

	upcomingEvents, err := d.repoEvent.GetFutureWithYear(ctx)
	if err != nil {
		return fmt.Errorf("get next events: %w", err)
	}

	dsas, err := d.repoDSA.GetByEvents(ctx, utils.SliceDereference(upcomingEvents))
	if err != nil {
		return err
	}

	var toCreate []*model.DSA

	for _, event := range upcomingEvents {
		act, activityOk := utils.SliceFind(activities, func(a activity) bool { return a.StartTime.Equal(event.StartTime) }) // TODO do not match based on start time as this can change
		dsa, dsaOk := utils.SliceFind(dsas, func(d *model.DSA) bool { return d.EventID == event.ID })

		switch {
		case activityOk:
			// Event is on the dsa website
			if !dsaOk {
				// Event has not been created yet locally
				toCreate = append(toCreate, &model.DSA{EventID: event.ID, DsaID: act.ID})
			} else if dsa.DsaID != 0 {
				// Both on the DSA website and locally in events, check if there is an update.
				updateBody := activityUpdate{}
				if act.Description != event.Description {
					updateBody.Description = event.Description
				}
				if !act.StartTime.Equal(event.StartTime) {
					updateBody.StartTime = event.StartTime
				}
				if !act.EndTime.Equal(event.EndTime) {
					updateBody.EndTime = event.EndTime
				}
				if act.Location != event.Location {
					updateBody.Location = event.Location
				}

				if (activityUpdate{}) != updateBody {
					if _, err := d.updateActivity(ctx, dsa.DsaID, updateBody); err != nil {
						return fmt.Errorf("update dsa activity: %w", err)
					}
				}

				if dsa.Deleted {
					// activity has been manually recreated on the dsa website
					dsa.Deleted = false
					if err := d.repoDSA.Update(ctx, *dsa); err != nil {
						return err
					}
				}
			}
		case !dsaOk:
			// Event is not on the dsa website yet and not created locally
			toCreate = append(toCreate, &model.DSA{EventID: event.ID})
		case dsa.DsaID != 0 && !dsa.Deleted:
			// This activity has been deleted on the DSA website, mark it as manually deleted.
			dsa.Deleted = true
			if err := d.repoDSA.Update(ctx, *dsa); err != nil {
				return fmt.Errorf("could not set dsa deleted to true: %w", err)
			}
		}
	}

	var errs []error

	for _, create := range toCreate {
		if err := d.repoDSA.Create(ctx, create); err != nil {
			errs = append(errs, err)
		}
	}

	if errs != nil {
		return errors.Join(errs...)
	}

	return nil
}

func (d *DSA) DeleteActivityByEvent(ctx context.Context, eventID int) error {
	dsa, err := d.repoDSA.GetByEventID(ctx, eventID)
	if err != nil {
		return fmt.Errorf("get dsa record by event: %w", err)
	}

	if dsa == nil {
		return fmt.Errorf("no dsa record having event id %d", eventID)
	}

	if dsa.DsaID != 0 {
		if _, err := d.deleteActivity(ctx, dsa.DsaID); err != nil {
			return fmt.Errorf("delete dsa activity: %w", err)
		}
	}

	return nil
}
