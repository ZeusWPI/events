package gitmate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type File struct {
	Name    string   `json:"name"`
	Path    string   `json:"path"`
	Type    FileType `json:"type"`
	Content []byte   `json:"content"`
}

type FileType string

const (
	TypeFile = "file"
	TypeDir  = "dir"
)

// Files will fetch all files for a filepath
// It's up to the user to use this function when a filepath to a directory is given
func (c *Client) Files(ctx context.Context, path string) ([]File, error) {
	req, err := c.newJSONRequest(ctx, "GET", "/contents/"+path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do http request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status for get files %s", resp.Status)
	}

	var files []File
	if err := decodeJSON(resp, &files); err != nil {
		return nil, fmt.Errorf("decode files: %w", err)
	}

	return files, nil
}

// File will fetch a file with it's contents
// It's up to the user to only use this function if a single file return is expected
func (c *Client) File(ctx context.Context, path string) ([]byte, error) {
	req, err := c.newJSONRequest(ctx, "GET", "/media/"+path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do http request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status for get file %s", resp.Status)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body bytes %w", err)
	}

	return bytes, nil
}

type FileCreate struct {
	Branch  string `json:"branch"`
	Content string `json:"content"`
	Message string `json:"message"`
}

func (c *Client) FileCreate(ctx context.Context, path string, body FileCreate) error {
	req, err := c.newJSONRequest(ctx, "POST", "/contents/"+path, body)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do http request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected HTTP status for file create: %s", resp.Status)
	}

	return nil
}

type Pull struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// Pulls returns all open pull requests (max 100)
func (c *Client) Pulls(ctx context.Context) ([]Pull, error) {
	req, err := c.newJSONRequest(ctx, "GET", "/pulls?state=open&page=1&limit=100", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do http request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status for get pulls: %s", resp.Status)
	}

	var pulls []Pull
	if err := decodeJSON(resp, &pulls); err != nil {
		return nil, fmt.Errorf("decode pulls: %w", err)
	}

	return pulls, nil
}

type PullCreate struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Base  string `json:"base"`
	Head  string `json:"head"`
}

// PullCreate creates a new pull request
func (c *Client) PullCreate(ctx context.Context, body PullCreate) error {
	req, err := c.newJSONRequest(ctx, "POST", "/pulls", body)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do http request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected HTTP status for pull create: %s", resp.Status)
	}

	return nil
}

type BranchCreate struct {
	Name string `json:"new_branch_name"`
	Base string `json:"old_branch_name"`
}

// BranchCreate creates a new branch

func (c *Client) BranchCreate(ctx context.Context, body BranchCreate) error {
	req, err := c.newJSONRequest(ctx, "POST", "/branches", body)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do http request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected HTTP status for branch create: %s", resp.Status)
	}

	return nil
}

func (c *Client) newJSONRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, fmt.Errorf("json encode: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url+path, &buf)
	if err != nil {
		return nil, fmt.Errorf("new http request %w", err)
	}

	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func decodeJSON[T any](resp *http.Response, out *T) error {
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(out)
}
