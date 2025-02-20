package website

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/ZeusWPI/events/internal/pkg/models"
	"github.com/ZeusWPI/events/pkg/util"
	"github.com/gocolly/colly"
	"go.uber.org/zap"
)

const eventURL = "https://zeus.gent/events"

// Get all event urls for a given academic year
func (w *Website) fetchEventURLSByAcademicYear(year models.AcademicYear) ([]string, error) {
	var urls []string
	var errs []error

	c := colly.NewCollector()
	c.OnHTML(".event-tile", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		parts := strings.Split(url, "/")
		if len(parts) != 5 {
			errs = append(errs, fmt.Errorf("Unable to get link of an event %s", url))
			return
		}
		urls = append(urls, parts[3])
	})

	err := c.Visit(fmt.Sprintf("%s/%s", eventURL, year.String()))
	if err != nil {
		return nil, err
	}

	c.Wait()

	if errs != nil {
		return nil, errors.Join(errs...)
	}

	return urls, nil
}

// UpdateEvent scrapes the website for event data and saves it
func (w *Website) UpdateEvent(event *models.Event) error {
	if event.URL == "" || (event.AcademicYear == models.AcademicYear{}) {
		return fmt.Errorf("Event has no URL or acdemic year: %+v", event)
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
			errs = append(errs, fmt.Errorf("Unable to parse event start time %s | %+v | %v", startRaw, *event, err))
		}
		event.StartTime = start
		// End time is not necessary
		endRaw := e.ChildAttr(".fa-ul > li:nth-child(1) > span:nth-child(4)", "content")
		end, err := time.Parse("2006-01-02T15:04:05-07:00", endRaw)
		if err == nil {
			event.EndTime = end
		}
	})

	err := c.Visit(fmt.Sprintf("%s/%s/%s", eventURL, event.AcademicYear.String(), event.URL))
	if err != nil {
		return err
	}

	c.Wait()

	if errs != nil {
		return errors.Join(errs...)
	}

	return w.eventRepo.Save(event)
}

// UpdateAllEvents synchronizes all events with the website
func (w *Website) UpdateAllEvents() error {
	zap.S().Debug("Updating all events")

	years, err := w.yearRepo.GetAll()
	if err != nil {
		return err
	}

	events, err := w.eventRepo.GetAll()
	if err != nil {
		return err
	}

	var errs []error
	var wg sync.WaitGroup
	for _, year := range years {
		wg.Add(1)
		go func(year models.AcademicYear) {
			defer wg.Done()

			urls, err := w.fetchEventURLSByAcademicYear(year)
			if err != nil {
				errs = append(errs, err)
				return
			}

			// Create / update each scraped event
			for _, url := range urls {
				var event *models.Event
				for _, e := range events {
					// Try to find existing event
					if e.URL == url && e.AcademicYear.Equal(year) {
						event = e
						break
					}
				}
				if event == nil {
					// Not found, create one
					event = &models.Event{URL: url, AcademicYear: year}
				}

				if err = w.UpdateEvent(event); err != nil {
					errs = append(errs, err)
				}
			}

			// Mark existing events that weren't found as deleted
			for _, event := range util.SliceFilter(events, func(e *models.Event) bool { return e.AcademicYear == year }) {
				if !slices.Contains(urls, event.URL) {
					if err = w.eventRepo.Delete(event); err != nil {
						errs = append(errs, err)
					}
				}
			}
		}(*year)
	}

	wg.Wait()

	return errors.Join(errs...)
}
