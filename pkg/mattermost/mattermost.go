// Package mattermost interacts with a mattermost instance
package mattermost

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ZeusWPI/events/pkg/config"
)

// Change this if you have a local mattermost instance
const url = "https://mattermost.zeus.gent/api/v4"

type Client struct {
	token string
}

func New() (*Client, error) {
	token := config.GetDefaultString("mattermost.token", "")
	if token == "" {
		return nil, errors.New("no mattermost token set")
	}

	return &Client{
		token: token,
	}, nil
}

type query struct {
	method string
	url    string
	body   io.Reader
	target any
}

func (c *Client) query(ctx context.Context, q query) error {
	req, err := http.NewRequestWithContext(ctx, q.method, fmt.Sprintf("%s/%s", url, q.url), q.body)
	if err != nil {
		return fmt.Errorf("new http request %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do http request %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected http status code %s", resp.Status)
	}

	if err = json.NewDecoder(resp.Body).Decode(q.target); err != nil {
		return fmt.Errorf("decode body to json %w", err)
	}

	return nil
}
