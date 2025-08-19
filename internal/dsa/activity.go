package dsa

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/pkg/utils"
)

const ActivitiesTask = "DSA activities update"
const CreateActivitesTask = "Create activites on DSA website"

type activityResponse struct {
	Page struct {
		Entries []activity `json:"entries"`
	} `json:"page"`
}

type activity struct {
	Association string    `json:"association"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
	ID          int       `json:"id"`
}

type activityCreate struct {
	Association string    `json:"association"`
	Description string    `json:"description"`
	EndTime     time.Time `json:"end_time"`
	StartTime   time.Time `json:"start_time"`
	Location    string    `json:"location"`
	Public      bool      `json:"public"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	Terrain     string    `json:"terrain"`
}

type activityUpdate struct {
	Association string    `json:"association,omitzero"`
	Description string    `json:"description,omitzero"`
	EndTime     time.Time `json:"end_time,omitzero"`
	StartTime   time.Time `json:"start_time,omitzero"`
	Location    string    `json:"location,omitzero"`
	Public      *bool     `json:"public,omitempty"`
	Title       string    `json:"title,omitzero"`
	Type        string    `json:"type,omitzero"`
	Terrain     string    `json:"terrain,omitzero"`
}

func (d *DSA) buildDsaURL(endpoint string, queries map[string]string) (string, error) {
	u, err := url.Parse(d.dsaURL)

	if err != nil {
		return "", err
	}

	u.Path, err = url.JoinPath(u.Path, endpoint)
	if err != nil {
		return "", err
	}

	query := url.Values{}

	for key, value := range queries {
		query.Set(key, value)
	}

	u.RawQuery = query.Encode()

	return u.String(), nil
}

func (d *DSA) doRequest(ctx context.Context, method string, url string, body any, target any) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, method, url, &buf)
	if err != nil {
		return fmt.Errorf("new http request %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", d.dsaKey)
	req.Header.Set("Content-Type", "application/json")

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

func (d *DSA) getActivities(ctx context.Context, target *activityResponse) error {
	dsaURL, err := d.buildDsaURL("activiteiten", map[string]string{
		"page_size":   "100",
		"association": d.abbreviation,
	})
	if err != nil {
		return err
	}

	if err = d.doRequest(ctx, http.MethodGet, dsaURL, nil, target); err != nil {
		return err
	}

	return nil
}

func (d *DSA) createActivity(ctx context.Context, body *activityCreate, target *activity) error {
	dsaURL, err := d.buildDsaURL("activiteiten", map[string]string{})
	if err != nil {
		return err
	}

	if err = d.doRequest(ctx, http.MethodPost, dsaURL, body, target); err != nil {
		return err
	}

	return nil
}

func (d *DSA) updateActivity(ctx context.Context, id int, body *activityUpdate, target *activity) error {
	dsaURL, err := d.buildDsaURL("activiteiten/"+strconv.Itoa(id), map[string]string{})
	if err != nil {
		return err
	}

	if err = d.doRequest(ctx, http.MethodPatch, dsaURL, body, target); err != nil {
		return err
	}

	return nil
}

func (d *DSA) deleteActivity(ctx context.Context, id int, target *activity) error {
	dsaURL, err := d.buildDsaURL("activiteiten/"+strconv.Itoa(id), map[string]string{})
	if err != nil {
		return err
	}

	if err = d.doRequest(ctx, http.MethodDelete, dsaURL, nil, target); err != nil {
		return err
	}
	return nil
}

func (d *DSA) CreateActivities(ctx context.Context) error {
	year, err := d.repoYear.GetLast(ctx)
	if err != nil {
		return err
	}

	if year == nil {
		return nil
	}

	var resp activityResponse
	if err := d.getActivities(ctx, &resp); err != nil {
		return err
	}

	activities := resp.Page.Entries

	events, err := d.repoEvent.GetByYearPopulated(ctx, year.ID)
	if err != nil {
		return err
	}

	now := time.Now()
	upcomingEvents := utils.SliceFilter(events, func(e *model.Event) bool { return e.StartTime.After(now) })

	dsas, err := d.repoDSA.GetByEvents(ctx, utils.SliceDereference(upcomingEvents))
	if err != nil {
		return err
	}

	uncreatedActivities := utils.SliceFilter(dsas, func(d *model.DSA) bool { return d.DsaID == 0 })

	for _, dsa := range uncreatedActivities {
		event, err := d.repoEvent.GetByID(ctx, dsa.EventID)
		act, activityOk := utils.SliceFind(activities, func(a activity) bool { return a.StartTime.Equal(event.StartTime) })

		if err != nil {
			return err
		}

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

			var response activity
			if err := d.createActivity(ctx, &activityCreate, &response); err != nil {
				return err
			}

			dsa.DsaID = response.ID

		} else {
			// Event is already on the dsa website, added manually
			dsa.DsaID = act.ID
		}

		if err := d.repoDSA.Update(ctx, dsa); err != nil {
			return err
		}
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

	for _, event := range upcomingEvents {
		act, activityOk := utils.SliceFind(activities, func(a activity) bool { return a.StartTime.Equal(event.StartTime) }) //TODO do not match based on start time as this can change
		dsa, dsaOk := utils.SliceFind(dsas, func(d *model.DSA) bool { return d.EventID == event.ID })

		if activityOk {
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
					var response activity
					d.updateActivity(ctx, dsa.DsaID, &updateBody, &response)
				}
			}
		} else if !dsaOk {
			// Event is not on the dsa website yet and not created locally
			toCreate = append(toCreate, &model.DSA{EventID: event.ID})
		} else if dsa.DsaID != 0 && dsa.Deleted == false {
			// This activity has been deleted on the DSA website, mark it as manually deleted.
			dsa.Deleted = true
			d.repoDSA.Update(ctx, dsa)
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
		return err
	}

	if dsa.EventID != 0 {
		var response activity
		if err := d.deleteActivity(ctx, dsa.DsaID, &response); err != nil {
			return err
		}
	}

	return nil
}
