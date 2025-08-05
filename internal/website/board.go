package website

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/pkg/utils"
)

const (
	BoardTask = "Board update"
	boardURL  = "https://raw.githubusercontent.com/ZeusWPI/zeus.ugent.be/master/data/bestuur.yaml"
)

type bestuurYAML struct {
	Data map[string][]struct {
		Role string `yaml:"rol"`
		Name string `yaml:"naam"`
	} `yaml:"data"`
}

func (w *Website) fetchAndParseBoard(ctx context.Context) ([]model.Board, error) {
	var raw bestuurYAML
	if err := w.github.FetchYaml(ctx, boardURL, &raw); err != nil {
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

func (w *Website) UpdateBoard(ctx context.Context) error {
	boards, err := w.fetchAndParseBoard(ctx)
	if err != nil {
		return err
	}

	years, err := w.yearRepo.GetAll(ctx)
	if err != nil {
		return nil
	}

	members, err := w.memberRepo.GetAll(ctx)
	if err != nil {
		return nil
	}

	oldBoards, err := w.boardRepo.GetAllPopulated(ctx)
	if err != nil {
		return err
	}

	var errs []error

	// Create or update
	for _, board := range boards {
		if exists := slices.ContainsFunc(oldBoards, func(b *model.Board) bool { return b.Equal(board) }); exists {
			// Exact copy already exists
			continue
		}

		// Is it an update?
		oldBoard, ok := utils.SliceFind(oldBoards, func(b *model.Board) bool { return b.EqualEntry(board) })
		if ok {
			// The board entry already exists with a different role or is_organizer value
			// This can either be because the website changed, the user logged in before the bestuur was updated or the user has 2 bestuurs roles
			// In the latter it will always update itself until the last role is reached
			// Update to preserve event assignments
			oldBoard.Role = board.Role
			oldBoard.IsOrganizer = board.IsOrganizer

			if err := w.boardRepo.Update(ctx, *oldBoard); err != nil {
				errs = append(errs, err)
			}

			continue
		}

		// Time to create

		// Get or create the member
		if member, ok := utils.SliceFind(members, func(m *model.Member) bool { return m.Equal(board.Member) }); ok {
			board.MemberID = member.ID
		} else {
			if err := w.memberRepo.Create(ctx, &board.Member); err != nil {
				errs = append(errs, err)
				continue
			}
			members = append(members, &board.Member)
			board.MemberID = board.Member.ID
		}

		// Get or create the year
		if year, ok := utils.SliceFind(years, func(y *model.Year) bool { return y.Equal(board.Year) }); ok {
			board.YearID = year.ID
		} else {
			if err := w.yearRepo.Create(ctx, &board.Year); err != nil {
				errs = append(errs, err)
				continue
			}
			years = append(years, &board.Year)
			board.YearID = board.Year.ID
		}

		if err := w.boardRepo.Create(ctx, &board); err != nil {
			errs = append(errs, err)
		}
	}

	// Delete old boards
	for _, board := range oldBoards {
		if !board.IsOrganizer {
			// Don't delete manually created board members
			// These entries will not be in boards
			continue
		}

		if exists := slices.ContainsFunc(boards, func(b model.Board) bool { return b.Equal(*board) }); !exists {
			if err := w.boardRepo.Delete(ctx, *board); err != nil {
				errs = append(errs, err)
			}
		}
	}

	if errs != nil {
		return errors.Join(errs...)
	}

	return nil
}
