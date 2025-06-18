package mattermost

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const apiURL = "api/v4"

type query struct {
	method string
	url    string
	body   io.Reader
	target any
}

func (m *Mattermost) query(ctx context.Context, q query) error {
	req, err := http.NewRequestWithContext(ctx, q.method, fmt.Sprintf("%s/%s/%s", m.url, apiURL, q.url), q.body)
	if err != nil {
		return fmt.Errorf("new http request %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do http request %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected http status code %s", resp.Status)
	}

	if err = json.NewDecoder(resp.Body).Decode(q.target); err != nil {
		return fmt.Errorf("decode body to json %w", err)
	}

	return nil
}
