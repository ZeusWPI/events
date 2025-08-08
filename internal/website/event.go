package website

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/pkg/utils"
	"gopkg.in/yaml.v3"
)

const (
	EventTask = "Event update"
	eventURL  = "https://api.github.com/repos/ZeusWPI/zeus.ugent.be/contents/content/events"
)

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
		parsed, err := time.Parse(layout, s)
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

func (c *Client) UpdateEvent(ctx context.Context) error {
	events, err := c.getEvents(ctx)
	if err != nil {
		return err
	}

	years, err := c.yearRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	oldEvents, err := c.eventRepo.GetAllWithYear(ctx)
	if err != nil {
		return err
	}

	var errs []error

	// Create / update events
	for _, event := range events {
		if exists := slices.ContainsFunc(oldEvents, func(e *model.Event) bool { return e.Equal(event) }); exists {
			continue
		}

		if year, ok := utils.SliceFind(years, func(y *model.Year) bool { return y.Equal(event.Year) }); ok {
			event.YearID = year.ID
		} else {
			if err := c.yearRepo.Create(ctx, &event.Year); err != nil {
				errs = append(errs, err)
				continue
			}
			event.YearID = event.Year.ID
		}

		var err error
		if exists := slices.ContainsFunc(oldEvents, func(e *model.Event) bool { return e.FileName == event.FileName }); exists {
			// Update
			err = c.eventRepo.Update(ctx, event)
		} else {
			// Create
			err = c.eventRepo.Create(ctx, &event)
		}

		if err != nil {
			errs = append(errs, err)
		}
	}

	// Refresh events
	newEventsP, err := c.eventRepo.GetAllWithYear(ctx)
	if err != nil {
		return err
	}
	newEvents := utils.SliceDereference(newEventsP)

	// Delete old events
	for _, event := range oldEvents {
		if exists := slices.ContainsFunc(newEvents, func(e model.Event) bool { return e.Equal(*event) }); !exists {
			if err := c.eventRepo.Delete(ctx, event.ID); err != nil {
				errs = append(errs, err)
			}
		}
	}

	if errs != nil {
		return errors.Join(errs...)
	}

	return nil
}
