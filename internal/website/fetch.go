package website

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gopkg.in/yaml.v3"
)

func (w *Website) fetchJSON(url string, target any) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("new http request %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+w.githubToken)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do http request %w", err)
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

func (w *Website) fetchMarkdown(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("new http request %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+w.githubToken)
	req.Header.Set("Accept", "application/vnd.github.raw")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("do http request %w", err)
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

func (w *Website) fetchYaml(url string, target any) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("new http request %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+w.githubToken)
	req.Header.Set("Accept", "application/vnd.github.raw")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do http request %w", err)
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
