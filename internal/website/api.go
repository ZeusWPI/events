package website

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
	"gopkg.in/yaml.v3"
)

// Board API

type bestuurYAML struct {
	Data map[string][]struct {
		Role string `yaml:"rol"`
		Name string `yaml:"naam"`
	} `yaml:"data"`
}

func (c *Client) fetchAndParseBoard(ctx context.Context) ([]model.Board, error) {
	var raw bestuurYAML
	if err := c.github.FetchYaml(ctx, boardURL, &raw); err != nil {
		return nil, err
	}

	var boards []model.Board
	for yearRange, members := range raw.Data {
		startYear, endYear, err := parseYearRange(yearRange)
		if err != nil {
			continue
		}

		for _, m := range members {
			boards = append(boards, model.Board{
				Role: m.Role,
				Member: model.Member{
					Name: m.Name,
				},
				Year: model.Year{
					Start: startYear,
					End:   endYear,
				},
				IsOrganizer: true,
			})
		}
	}

	return boards, nil
}

func parseYearRange(s string) (int, int, error) {
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid year range %s", s)
	}

	startSuffix, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}

	endSuffix, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}

	// This will break in 2090
	startYear := 1900 + startSuffix
	if startSuffix < 90 {
		startYear = 2000 + startSuffix
	}
	endYear := 1900 + endSuffix
	if endSuffix < 90 {
		endYear = 2000 + endSuffix
	}

	return startYear, endYear, nil
}

// Events API

var headerRegex = regexp.MustCompile(`(?s)^---[^\n]*\n(.*?)\n---`)

type fileMeta struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	DownloadURL string `json:"download_url"`
}

type header struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Time        string `yaml:"time"`
	End         string `yaml:"end"`
	Location    string `yaml:"location"`
}

func parseYearDir(name string) (int, int, error) {
	parts := strings.Split(name, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid year dir: %s", name)
	}

	start, _ := parseYearSuffix(parts[0])
	end, _ := parseYearSuffix(parts[1])

	return start, end, nil
}

func parseYearSuffix(s string) (int, error) {
	parsedYear, err := time.Parse("06", s)
	if err != nil {
		return 0, err
	}

	year := parsedYear.Year()
	if year < 1990 {
		year += 100
	}
	return year, nil
}

func parseTime(s string) (time.Time, error) {
	location, err := time.LoadLocation("Europe/Brussels")
	if err != nil {
		return time.Time{}, fmt.Errorf("load brussels time zone %w", err)
	}

	// Zeus can't be consistent in the date formatting...
	layouts := []string{
		"02-01-2006 15:04",
		"2-01-2006 15:04",
		"02-1-2006 15:04",
		"2-1-2006 15:04",
		"2006-01-02 15:04",
		"2006-1-02 15:04",
		"2006-01-2 15:04",
		"2006-1-2 15:04",
		"02-01-2006 15h",
		"2-01-2006 15h",
		"02-1-2006 15h",
		"2-1-2006 15h",
		"02-01-2006 15h04",
		"2-01-2006 15h04",
		"02-1-2006 15h04",
		"2-1-2006 15h04",
		"2006-01-02T15:04:05-07:00",
		"2006-1-02T15:04:05-07:00",
		"2006-01-2T15:04:05-07:00",
		"2006-1-2T15:04:05-07:00",
		"02-01-2006",
		"2-01-2006",
		"02-1-2006",
		"2-1-2006",
	}

	if ok := strings.Contains(s, "24:00"); ok {
		s = strings.Replace(s, "24:00", "23:59", 1)
	}

	for _, layout := range layouts {
		parsed, err := time.ParseInLocation(layout, s, location)
		if err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("parse time %s", s)
}

func (c *Client) parseEventFile(ctx context.Context, dirName string, f fileMeta) (model.Event, error) {
	if !strings.HasSuffix(f.Name, ".md") {
		return model.Event{}, fmt.Errorf("invalid file %+v", f)
	}

	mdContent, err := c.github.FetchMarkdown(ctx, f.DownloadURL)
	if err != nil {
		return model.Event{}, err
	}

	match := headerRegex.FindStringSubmatch(mdContent)
	if len(match) < 2 {
		return model.Event{}, fmt.Errorf("no header found in %s", f.Path)
	}

	var head header
	if err := yaml.Unmarshal([]byte(match[1]), &head); err != nil {
		return model.Event{}, err
	}

	startYear, endYear, err := parseYearDir(dirName)
	if err != nil {
		return model.Event{}, err
	}

	startTime, err := parseTime(head.Time)
	if err != nil {
		return model.Event{}, fmt.Errorf("invalid start time in %s: %w", f.Path, err)
	}

	endTime := time.Time{}
	if head.End != "" {
		endTime, err = parseTime(head.End)
		if err != nil {
			return model.Event{}, fmt.Errorf("invalid end time in %s: %w", f.Path, err)
		}
	}

	return model.Event{
		FileName:    strings.TrimSuffix(f.Name, ".md"),
		Name:        head.Title,
		Description: head.Description,
		StartTime:   startTime,
		EndTime:     endTime,
		Location:    head.Location,
		Year: model.Year{
			Start: startYear,
			End:   endYear,
		},
	}, nil
}

func (c *Client) getEvents(ctx context.Context) ([]model.Event, error) {
	var yearDirs []fileMeta
	if err := c.github.FetchJSON(ctx, eventURL, &yearDirs); err != nil {
		return nil, fmt.Errorf("fetch year dirs: %w", err)
	}

	var all []model.Event
	var errs []error

	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, dir := range yearDirs {
		if dir.Type != "dir" {
			continue
		}

		wg.Add(1)

		go func(dir fileMeta) {
			defer wg.Done()

			var files []fileMeta
			url := fmt.Sprintf("%s/%s", eventURL, dir.Name)

			if err := c.github.FetchJSON(ctx, url, &files); err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("failed to fetch files for %s: %w", dir.Name, err))
				mu.Unlock()

				return
			}

			for _, file := range files {
				if file.Type != "file" || !strings.HasSuffix(file.Name, ".md") {
					continue
				}

				event, err := c.parseEventFile(ctx, dir.Name, file)
				if err != nil {
					mu.Lock()
					errs = append(errs, err)
					mu.Unlock()

					return
				}

				mu.Lock()
				all = append(all, event)
				mu.Unlock()
			}
		}(dir)
	}

	wg.Wait()

	if len(errs) != 0 {
		return nil, errors.Join(errs...)
	}

	return all, nil
}
