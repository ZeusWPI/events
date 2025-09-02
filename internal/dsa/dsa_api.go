package dsa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

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
		return "", fmt.Errorf("dsaURL could not be parsed: %w", err)
	}

	u.Path, err = url.JoinPath(u.Path, endpoint)
	if err != nil {
		return "", fmt.Errorf("could not add endpoint to path: %w", err)
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
		return fmt.Errorf("could not encode the body struct: %w", err)
	}

	if d.development && method != "GET" {
		// Do not do the actual request in development
		zap.S().Infof("Mock request: %s %s\n\tBody: %+v", method, url, body)
		return nil
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

func (d *DSA) getActivities(ctx context.Context) ([]activity, error) {
	var response activityResponse
	dsaURL, err := d.buildDsaURL("activiteiten", map[string]string{
		"page_size":   "100",
		"association": d.abbreviation,
	})
	if err != nil {
		return nil, fmt.Errorf("build dsa url: %w", err)
	}

	if err = d.doRequest(ctx, http.MethodGet, dsaURL, nil, &response); err != nil {
		return nil, fmt.Errorf("do dsa request: %w", err)
	}

	return response.Page.Entries, nil
}

func (d *DSA) createActivity(ctx context.Context, body activityCreate) (activity, error) {
	var response activity
	dsaURL, err := d.buildDsaURL("activiteiten", map[string]string{})
	if err != nil {
		return response, fmt.Errorf("build dsa url: %w", err)
	}

	if err = d.doRequest(ctx, http.MethodPost, dsaURL, body, &response); err != nil {
		return response, fmt.Errorf("do dsa request: %w", err)
	}

	return response, nil
}

func (d *DSA) updateActivity(ctx context.Context, id int, body activityUpdate) (activity, error) {
	var response activity
	dsaURL, err := d.buildDsaURL("activiteiten/"+strconv.Itoa(id), map[string]string{})
	if err != nil {
		return response, fmt.Errorf("build dsa url: %w", err)
	}

	if err = d.doRequest(ctx, http.MethodPatch, dsaURL, body, &response); err != nil {
		return response, fmt.Errorf("do dsa request: %w", err)
	}

	return response, nil
}

func (d *DSA) deleteActivity(ctx context.Context, id int) (activity, error) {
	var response activity
	dsaURL, err := d.buildDsaURL("activiteiten/"+strconv.Itoa(id), map[string]string{})
	if err != nil {
		return response, fmt.Errorf("build dsa url: %w", err)
	}

	if err = d.doRequest(ctx, http.MethodDelete, dsaURL, nil, &response); err != nil {
		return response, fmt.Errorf("do dsa request: %w", err)
	}
	return response, nil
}
