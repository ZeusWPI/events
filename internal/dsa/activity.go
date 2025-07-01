package dsa

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
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

func (d *DSA) BuildDsaURL(endpoint string, queries map[string]string) (string, error) {
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

func (d *DSA) getActivities(ctx context.Context, target *activityResponse) error {
	dsaURL, err := d.BuildDsaURL("activiteiten", map[string]string{
		"page_size":   "100",
		"association": d.abbreviation,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", dsaURL, nil)
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

func (d *DSA) createActivity(ctx context.Context, body *activityCreate, target *activity) error {
	dsaURL, err := d.BuildDsaURL("activiteiten", map[string]string{})
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", dsaURL, &buf)
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

		if !activityOk {
			// Event is not yet on the dsa website

			if err != nil {
				return err
			}

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
		activity, activityOk := utils.SliceFind(activities, func(a activity) bool { return a.StartTime.Equal(event.StartTime) })
		_, dsaOk := utils.SliceFind(dsas, func(d *model.DSA) bool { return d.EventID == event.ID })

		if activityOk {
			// Event is on the dsa website
			if !dsaOk {
				toCreate = append(toCreate, &model.DSA{EventID: event.ID, DsaID: activity.ID})
			}
		} else if !dsaOk {
			// Event is not on the dsa website yet (or has been removed)
			toCreate = append(toCreate, &model.DSA{EventID: event.ID})
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
