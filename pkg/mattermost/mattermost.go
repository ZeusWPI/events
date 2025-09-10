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

const apiURL = "api/v4"

type Client struct {
	token string
	url   string
}

func New() (*Client, error) {
	token := config.GetDefaultString("mattermost.token", "")
	if token == "" {
		return nil, errors.New("no mattermost token set")
	}

	url := config.GetDefaultString("mattermost.url", "")
	if url == "" {
		return nil, errors.New("no mattermost url set")
	}

	return &Client{
		token: token,
		url:   url,
	}, nil
}

type query struct {
	method string
	url    string
	body   io.Reader
	target any
}

func (c *Client) query(ctx context.Context, q query) error {
	req, err := http.NewRequestWithContext(ctx, q.method, fmt.Sprintf("%s/%s/%s", c.url, apiURL, q.url), q.body)
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
