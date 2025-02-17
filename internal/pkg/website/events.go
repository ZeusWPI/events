package website

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ZeusWPI/events/internal/pkg/models"
	"github.com/ZeusWPI/events/pkg/util"
	"github.com/gocolly/colly"
	"go.uber.org/zap"
)

// Warning: Webscraping results in ugly code

const eventURL = "https://zeus.gent/events"

// Get all academic years
func (w *Website) fetchAllAcademicYears() ([]string, error) {
	zap.S().Debug("Fetching academic years")

	var years []string
	var errs []error

	c := colly.NewCollector()
	c.OnHTML(".menu-list", func(e *colly.HTMLElement) {
		yearsRaw := e.ChildAttrs("a", "href")
		if len(yearsRaw) < 2 {
			// Will only happen if someone nukes the Zeus WPI website
			return
		}
		sort.Strings(yearsRaw)
		// The current year (represented by '#') is now the first element and last year is the last element.
		lastYear, err := getAcademicYear(yearsRaw[len(yearsRaw)-1])
		if err != nil {
			errs = append(errs, err)
			return
		}

		currentYear, err := incrementYear(lastYear)
		if err != nil {
			errs = append(errs, err)
			return
		}

		for _, year := range yearsRaw[1:] {
			y, err := getAcademicYear(year)
			if err != nil {
				errs = append(errs, err)
				continue
			}

			years = append(years, y)
		}
		years = append(years, currentYear)
	})

	err := c.Visit(eventURL)
	if err != nil {
		return nil, fmt.Errorf("Unable to visit Zeus WPI website %s | %w", eventURL, err)
	}

	c.Wait()

	if errs != nil {
		return nil, errors.Join(errs...)
	}

	return years, nil
}

// Get all event urls for a given academic year
func (w *Website) fetchEventURLSByAcademicYear(year string) ([]string, error) {
	zap.S().Debug("Fetching event URLS for academic year ", year)

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

	err := c.Visit(fmt.Sprintf("%s/%s", eventURL, year))
	if err != nil {
		return nil, err
	}

	c.Wait()

	if errs != nil {
		return nil, errors.Join(errs...)
	}

	return urls, nil
}

// Update scrapes the website for event data and saves it
func (w *Website) Update(event *models.Event) error {
	zap.S().Debug("Updating event ", event.Name, event.URL, event.AcademicYear)

	if event.URL == "" || event.AcademicYear == "" {
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
			errs = append(errs, fmt.Errorf("Unable to parse event start time %s | %+v | %w", startRaw, *event, err))
		}
		event.StartTime = start
		// End time is not necessary
		endRaw := e.ChildAttr(".fa-ul > li:nth-child(1) > span:nth-child(4)", "content")
		end, err := time.Parse("2006-01-02T15:04:05-07:00", endRaw)
		if err == nil {
			event.EndTime = end
		}
	})

	err := c.Visit(fmt.Sprintf("%s/%s/%s", eventURL, event.AcademicYear, event.URL))
	if err != nil {
		return err
	}

	c.Wait()

	if errs != nil {
		return errors.Join(errs...)
	}

	return w.eventRepo.Save(event)
}

// UpdateAll synchronizes all events with the website
func (w *Website) UpdateAll() error {
	zap.S().Debug("Updating all events")

	years, err := w.fetchAllAcademicYears()
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
		go func(year string) {
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
					if e.URL == url && e.AcademicYear == year {
						event = e
						continue
					}
				}
				if event == nil {
					// Not found, create one
					event = &models.Event{URL: url, AcademicYear: year}
				}

				if err = w.Update(event); err != nil {
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
		}(year)
	}

	wg.Wait()

	return errors.Join(errs...)
}
