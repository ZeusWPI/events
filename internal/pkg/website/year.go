package website

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/gocolly/colly"
)

const (
	// YearTask is the name of the recurring task that updates all years
	YearTask = "Years Update"
	yearURL  = "https://zeus.gent/events"
)

// Get all years
func (w *Website) fetchAllYears() ([]string, error) {
	var years []string
	var errs []error

	c := colly.NewCollector()
	c.OnHTML(".menu-list", func(e *colly.HTMLElement) {
		yearsRaw := e.ChildAttrs("a", "href")
		if len(yearsRaw) < 2 {
			// Will only happen if someone nukes the Zeus WPI website
			return
		}
		sort.Sort(sort.Reverse(sort.StringSlice(yearsRaw)))
		// The current year (represented by '#') is now the last element and last year is the first element.
		lastYear, err := getYear(yearsRaw[0])
		if err != nil {
			errs = append(errs, err)
			return
		}

		currentYear, err := incrementYear(lastYear)
		if err != nil {
			errs = append(errs, err)
			return
		}
		years = append(years, currentYear)

		for _, year := range yearsRaw[:len(yearsRaw)-1] {
			y, err := getYear(year)
			if err != nil {
				errs = append(errs, err)
				continue
			}

			years = append(years, y)
		}
	})

	err := c.Visit(yearURL)
	if err != nil {
		return nil, fmt.Errorf("unable to visit Zeus WPI website %s | %w", yearURL, err)
	}

	c.Wait()

	if errs != nil {
		return nil, errors.Join(errs...)
	}

	return years, nil
}

// UpdateAllYears years
func (w *Website) UpdateAllYears() error {
	yearsWebsite, err := w.fetchAllYears()
	if err != nil {
		return err
	}

	yearsDB, err := w.yearRepo.GetAll(context.Background())
	if err != nil {
		return err
	}

	// Don't delete because of foreign key nightmare with events
	var yearsToAdd []string

	// Both are sorted in the same order
	for i := 0; i < len(yearsWebsite) && i < len(yearsDB); i++ {
		if yearsWebsite[i] != yearsDB[i].String() {
			yearsToAdd = append(yearsToAdd, yearsWebsite[i])
		}
	}

	for i := len(yearsDB); i < len(yearsWebsite); i++ {
		yearsToAdd = append(yearsToAdd, yearsWebsite[i])
	}

	var errs []error
	for _, y := range yearsToAdd {
		parts := strings.Split(y, "-")
		start, err1 := strconv.Atoi("20" + parts[0]) // Come find me when this breaks
		end, err2 := strconv.Atoi("20" + parts[1])
		if err1 != nil || err2 != nil {
			errs = append(errs, fmt.Errorf("unable to convert string year to int %s | %w | %w", y, err1, err2))
		}

		if err := w.yearRepo.Save(context.Background(), &model.Year{
			StartYear: start, EndYear: end,
		}); err != nil {
			errs = append(errs, err)
		}
	}

	if errs != nil {
		return fmt.Errorf("unable to update all years %w", errors.Join(errs...))
	}

	return nil
}
