package website

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ZeusWPI/events/internal/pkg/model"
	"github.com/gocolly/colly"
	"go.uber.org/zap"
)

const boardURL = "https://zeus.gent/about/oud-bestuur/"

func (w *Website) fetchAllBoards() ([]model.Board, error) {
	var boards []model.Board
	var errs []error

	yearRegex := regexp.MustCompile(`(\d{4})\s*–\s*(\d{4})`)

	c := colly.NewCollector()
	c.OnHTML("h2", func(e *colly.HTMLElement) {
		match := yearRegex.FindStringSubmatch(e.Text)
		if len(match) != 3 {
			return
		}

		yearStart, err1 := strconv.Atoi(match[1])
		yearEnd, err2 := strconv.Atoi(match[2])

		if err1 != nil || err2 != nil {
			return
		}

		e.DOM.Next().Find("tr").Each(func(_ int, row *goquery.Selection) {
			cells := row.Find("td")
			if len(cells.Nodes) != 2 {
				return
			}

			role := strings.TrimSpace(cells.Eq(0).Text())
			name := strings.TrimSpace(cells.Eq(1).Text())
			boards = append(boards, model.Board{
				Member: model.Member{
					Name: name,
				},
				AcademicYear: model.AcademicYear{
					StartYear: yearStart,
					EndYear:   yearEnd,
				},
				Role: role,
			})
		})
	})

	err := c.Visit(boardURL)
	if err != nil {
		return nil, fmt.Errorf("Unable to visit Zeus WPI website %s | %v", boardURL, err)
	}

	c.Wait()

	if errs != nil {
		return nil, errors.Join(errs...)
	}

	return boards, nil
}

// UpdateAllBoards updates all boards
func (w *Website) UpdateAllBoards() error {
	zap.S().Debug("Updating entire board")

	// Fetch all data from website or DB
	boardsWebsite, err := w.fetchAllBoards()
	if err != nil {
		return err
	}

	years, err := w.yearRepo.GetAll()
	if err != nil {
		return nil
	}

	members, err := w.memberRepo.GetAll()
	if err != nil {
		return nil
	}

	boards, err := w.boardRepo.GetAll()
	if err != nil {
		return err
	}

	var errs []error
	for _, board := range boardsWebsite {
		// Look for new boards
		exists := slices.ContainsFunc(boards, func(b *model.Board) bool {
			return b.Equal(board)
		})

		if !exists {
			// Find existing member
			newMember := true
			for _, member := range members {
				if member.Equal(board.Member) {
					board.Member = *member
					newMember = false
					break
				}
			}

			// Find existing academic year
			newYear := true
			for _, year := range years {
				if year.Equal(board.AcademicYear) {
					newYear = false
					board.AcademicYear = *year
					break
				}
			}

			err := w.boardRepo.Save(&board)
			if err != nil {
				errs = append(errs, err)
				break
			}

			// Update the existing member and year list
			if newMember {
				members = append(members, &board.Member)
			}
			if newYear {
				years = append(years, &board.AcademicYear)
			}
		}
	}

	if errs != nil {
		return errors.Join(errs...)
	}

	return nil
}
