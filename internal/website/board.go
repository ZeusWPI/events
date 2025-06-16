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

func (w *Website) fetchAndParseBoard() ([]model.Board, error) {
	var raw bestuurYAML
	if err := w.fetchYaml(boardURL, &raw); err != nil {
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
	boards, err := w.fetchAndParseBoard()
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

	oldBoards, err := w.boardRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	var errs []error

	// Create new boards
	for _, board := range boards {
		if exists := slices.ContainsFunc(oldBoards, func(b *model.Board) bool { return b.Equal(board) }); exists {
			continue
		}

		if member, ok := utils.SliceFind(members, func(m *model.Member) bool { return m.Equal(board.Member) }); ok {
			board.MemberID = member.ID
		} else {
			if err := w.memberRepo.Create(ctx, &board.Member); err != nil {
				errs = append(errs, err)
				continue
			}
			board.MemberID = board.Member.ID
		}

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
