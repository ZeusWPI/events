// Package github interacts with the github API
package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/ZeusWPI/events/pkg/config"
	"gopkg.in/yaml.v3"
)

type Client struct {
	token string

	mu        sync.Mutex
	remaining int
	reset     time.Time
}

func New() (*Client, error) {
	token := config.GetDefaultString("github.token", "")
	if token == "" {
		return nil, errors.New("no github token set")
	}

	return &Client{
		token:     token,
		remaining: 1,
		reset:     time.Now(),
	}, nil
}

func (c *Client) doRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	// Rate limit check
	c.mu.Lock()

	if time.Now().Before(c.reset) && c.remaining <= 0 {
		sleepFor := time.Until(c.reset) + time.Second
		c.mu.Unlock()

		select {
		case <-time.After(sleepFor):
		case <-ctx.Done():
			return nil, ctx.Err()
		}

		c.mu.Lock() // Avoid the runtime error by unlocking afterwards
	}
	c.mu.Unlock()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do http request %w", err)
	}

	// Update rate limit info
	if limit := resp.Header.Get("X-RateLimit-Remaining"); limit != "" {
		if rem, err := strconv.Atoi(limit); err == nil {
			c.mu.Lock()
			c.remaining = rem
			c.mu.Unlock()
		}
	}

	if reset := resp.Header.Get("X-RateLimit-Reset"); reset != "" {
		if ts, err := strconv.ParseInt(reset, 10, 64); err == nil {
			c.mu.Lock()
			c.reset = time.Unix(ts, 0)
			c.mu.Unlock()
		}
	}

	return resp, nil
}

func (c *Client) FetchJSON(ctx context.Context, url string, target any) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("new http request %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status %s", resp.Status)
	}

	if err = json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("decode body to json %w", err)
	}

	return nil
}

func (c *Client) FetchMarkdown(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("new http request %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github.raw")

	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected HTTP status %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading body %w", err)
	}

	return string(body), nil
}

func (c *Client) FetchYaml(ctx context.Context, url string, target any) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("new http request %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github.raw")

	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status %s", resp.Status)
	}

	if err = yaml.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("decode body to yaml %w", err)
	}

	return nil
}
