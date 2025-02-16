// Package website scrapes the Zeus WPI website to get all event data
package website

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ZeusWPI/events/internal/pkg/db/repository"
	"github.com/ZeusWPI/events/internal/pkg/models"
	"github.com/gocolly/colly"
)

const baseURL = "https://zeus.gent"

// Website represents the ZeusWPI website and all related functions
type Website struct {
	EventRepo repository.Event
}

// New creates a new website instance
func New(repo repository.Repository) *Website {
	return &Website{
		EventRepo: repo.NewEvent(),
	}
}

// UpdateAll fetches all events and related data
func (w *Website) UpdateAll() error {
	// Fetch all events
	var events []*models.Event
	errs := make([]error, 0)

	c := colly.NewCollector()
	c.OnHTML(".event-tile", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if link == "" {
			errs = append(errs, errors.New("Unable to get link of an event"))
			return
		}

		event := models.Event{}
		event.URL = link
		events = append(events, &event)
	})

	err := c.Visit(fmt.Sprintf("%s/events", baseURL))
	if err != nil {
		return err
	}

	c.Wait()

	if errs != nil {
		return errors.Join(errs...)
	}

	for _, e := range events {
		errs = append(errs, w.Update(e))
	}

	return errors.Join(errs...)
}

// Update scrapes the website for new data for a given event
func (w *Website) Update(event *models.Event) error {
	if event.URL == "" {
		return fmt.Errorf("Event has no URL: %+v", event)
	}

	var errs []error

	c := colly.NewCollector()
	c.OnHTML(".header-text", func(e *colly.HTMLElement) {
		event.Name = e.ChildText(".is-1-responsive > b:nth-child(1)")
		event.Description = e.ChildText(".is-3-responsive")

		startRaw := e.ChildAttr(".fa-ul > li:nth-child(1) > span:nth-child(2)", "content")
		start, err := time.Parse("2006-01-02T15:04:05-07:00", startRaw)
		if err != nil {
			errs = append(errs, fmt.Errorf("Unable to parse event %+v start time %s", *event, startRaw))
		}
		event.StartTime = start
		endRaw := e.ChildAttr(".fa-ul > li:nth-child(1) > span:nth-child(4)", "content")
		end, err := time.Parse("2006-01-02T15:04:05-07:00", endRaw)
		if err != nil {
			errs = append(errs, fmt.Errorf("Unable to parse event %+v end time %s", *event, startRaw))
		}
		event.EndTime = end

		urlParts := strings.Split(event.URL, "/")
		if len(urlParts) != 5 {
			errs = append(errs, fmt.Errorf("Unable to parse URL to retrieve academic year %+v", *event))
		} else {
			event.AcademicYear = urlParts[2]
		}

		event.Location = e.ChildText(".fa-ul > li:nth-child(2) > span:nth-child(2)")
	})

	err := c.Visit(fmt.Sprintf("%s%s", baseURL, event.URL))
	if err != nil {
		return err
	}

	c.Wait()

	if errs != nil {
		return errors.Join(errs...)
	}

	return w.EventRepo.Save(event)
}
