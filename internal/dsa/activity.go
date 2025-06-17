package dsa

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/pkg/utils"
)

const ActivitiesTask = "DSA activities update"

type activityResponse struct {
	Page struct {
		Entries []activity `json:"entries"`
	} `json:"page"`
}

type activity struct {
	Association string    `json:"association"`
	StartTime   time.Time `json:"start_time"`
}

func (d *DSA) getActivities(ctx context.Context, target *activityResponse) error {
	req, err := http.NewRequestWithContext(ctx, "GET", d.dsaURL, nil)
	if err != nil {
		return fmt.Errorf("new http request %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", d.dsaKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do http request %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected http status code %s", resp.Status)
	}

	if err = json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("decode body to json %w", err)
	}

	return nil
}

func (d *DSA) UpdateActivities(ctx context.Context) error {
	var resp activityResponse
	if err := d.getActivities(ctx, &resp); err != nil {
		return err
	}
	activities := resp.Page.Entries

	year, err := d.repoYear.GetLast(ctx)
	if err != nil {
		return err
	}
	if year == nil {
		return nil
	}

	events, err := d.repoEvent.GetByYearPopulated(ctx, year.ID)
	if err != nil {
		return err
	}

	// DSA api only shows upcoming events
	now := time.Now()
	upcomingEvents := utils.SliceFilter(events, func(e *model.Event) bool { return e.StartTime.After(now) })

	dsas, err := d.repoDSA.GetByEvents(ctx, utils.SliceDereference(upcomingEvents))
	if err != nil {
		return err
	}

	var toCreate []*model.DSA
	var toDelete []int

	for _, event := range upcomingEvents {
		_, activityOk := utils.SliceFind(activities, func(a activity) bool { return a.StartTime.Equal(event.StartTime) })
		dsa, dsaOk := utils.SliceFind(dsas, func(d *model.DSA) bool { return d.EventID == event.ID })

		if activityOk {
			// Event is on the dsa website
			if !dsaOk {
				toCreate = append(toCreate, &model.DSA{EventID: event.ID, Entry: true})
			}
		} else {
			// Event is not on the dsa website yet (or has been removed)
			if dsaOk {
				toDelete = append(toDelete, dsa.ID)
			}
		}
	}

	var errs []error

	for _, create := range toCreate {
		if err := d.repoDSA.Create(ctx, create); err != nil {
			errs = append(errs, err)
		}
	}

	for _, delete := range toDelete {
		if err := d.repoDSA.Delete(ctx, delete); err != nil {
			errs = append(errs, err)
		}
	}

	if errs != nil {
		return errors.Join(errs...)
	}

	return nil
}
