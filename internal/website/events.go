package website

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/pkg/util"
	"github.com/gocolly/colly"
)

const (
	// EventTask is the name of the recurring task that updates all events
	EventTask      = "Events Update"
	eventURL       = "https://zeus.gent/events"
	eventStartYear = 2000
)

// Get all event urls for a given year
func (w *Website) fetchEventURLSByYear(year model.Year) ([]string, error) {
	if year.StartYear < eventStartYear {
		return nil, nil
	}

	var urls []string
	var errs []error

	c := colly.NewCollector()
	c.OnHTML(".event-tile", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		parts := strings.Split(url, "/")
		if len(parts) != 5 {
			errs = append(errs, fmt.Errorf("unable to get link of an event %s", url))
			return
		}
		urls = append(urls, parts[3])
	})

	url := fmt.Sprintf("%s/%s", eventURL, year.String())
	err := c.Visit(url)
	if err != nil {
		return nil, fmt.Errorf("unable to visit url %s | %w", url, err)
	}

	c.Wait()

	if errs != nil {
		return nil, errors.Join(errs...)
	}

	return urls, nil
}

// UpdateEvent scrapes the website for event data and saves it
func (w *Website) UpdateEvent(event *model.Event) error {
	if event.URL == "" || (event.Year == model.Year{}) {
		return fmt.Errorf("event has no URL or acdemic year: %+v", event)
	}

	var errs []error

	c := colly.NewCollector()
	c.OnHTML(".header-text", func(e *colly.HTMLElement) {
		event.Name = e.ChildText(".is-1-responsive > b:nth-child(1)")
		event.Description = e.ChildText(".is-3-responsive")
		event.Location = e.ChildText(".fa-ul > li:nth-child(2) > span:nth-child(2)")

		startRaw := e.ChildAttr(".fa-ul > li:nth-child(1) > span:nth-child(2)", "content")
		start, err := time.Parse("2006-01-02T15:04:05-07:00", startRaw)
		if err != nil {
			errs = append(errs, fmt.Errorf("unable to parse event start time %s | %+v | %w", startRaw, *event, err))
		}
		event.StartTime = start
		// End time is not necessary
		endRaw := e.ChildAttr(".fa-ul > li:nth-child(1) > span:nth-child(4)", "content")
		end, err := time.Parse("2006-01-02T15:04:05-07:00", endRaw)
		if err == nil {
			event.EndTime = end
		}
	})

	url := fmt.Sprintf("%s/%s/%s", eventURL, event.Year.String(), event.URL)
	err := c.Visit(url)
	if err != nil {
		return fmt.Errorf("unable to visit url %s | %w", url, err)
	}

	c.Wait()

	if errs != nil {
		return errors.Join(errs...)
	}

	return w.eventRepo.Save(context.Background(), event)
}

// UpdateAllEvents synchronizes all events with the website
func (w *Website) UpdateAllEvents() error {
	years, err := w.yearRepo.GetAll(context.Background())
	if err != nil {
		return err
	}

	events, err := w.eventRepo.GetAllWithYear(context.Background())
	if err != nil {
		return err
	}

	var errs []error
	var wg sync.WaitGroup
	for _, year := range years {
		wg.Add(1)
		go func(year model.Year) {
			defer wg.Done()

			urls, err := w.fetchEventURLSByYear(year)
			if err != nil {
				errs = append(errs, err)
				return
			}

			// Create / update each scraped event
			for _, url := range urls {
				var event *model.Event
				for _, e := range events {
					// Try to find existing event
					if e.URL == url && e.Year.Equal(year) {
						event = e
						break
					}
				}
				if event == nil {
					// Not found, create one
					event = &model.Event{URL: url, Year: year}
				}

				if err = w.UpdateEvent(event); err != nil {
					errs = append(errs, err)
				}
			}

			// Mark existing events that weren't found as deleted
			for _, event := range util.SliceFilter(events, func(e *model.Event) bool { return e.Year == year }) {
				if !slices.Contains(urls, event.URL) {
					if err = w.eventRepo.Delete(context.Background(), event); err != nil {
						errs = append(errs, err)
					}
				}
			}
		}(*year)
	}

	wg.Wait()

	return errors.Join(errs...)
}
