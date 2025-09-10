package website

import (
	"context"
	"fmt"
	"slices"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/pkg/utils"
)

const boardURL = "https://raw.githubusercontent.com/ZeusWPI/zeus.ugent.be/master/data/bestuur.yaml"

func (c *Client) SyncBoard(ctx context.Context) error {
	websiteBoards, err := c.fetchAndParseBoard(ctx)
	if err != nil {
		return err
	}

	dbYears, err := c.yearRepo.GetAll(ctx)
	if err != nil {
		return nil
	}

	dbMembers, err := c.memberRepo.GetAll(ctx)
	if err != nil {
		return nil
	}

	dbBoards, err := c.boardRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	// Create or update
	for _, board := range websiteBoards {
		if exists := slices.ContainsFunc(dbBoards, func(b *model.Board) bool { return b.Equal(board) }); exists {
			// Exact copy already exists
			continue
		}

		// Is it an update?
		if oldBoard, ok := utils.SliceFind(dbBoards, func(b *model.Board) bool { return b.EqualEntry(board) }); ok {
			// Both the website and the local database contain this board member
			// but they differ slightly, let's bring it up to date.
			// This situation can happen if
			//   - The website entry changed (e.g. a new role)
			//   - The user logged in before the board sync happened
			//   - The user has multiple role assignments (happened in the beginning of Zeus WPI)
			board.ID = oldBoard.ID
			board.MemberID = oldBoard.MemberID
			board.YearID = oldBoard.YearID

			if err := c.boardRepo.Update(ctx, board); err != nil {
				return fmt.Errorf("updating board entry for old board %+v | %w", *oldBoard, err)
			}

			continue
		}

		// We now know it's a new board member
		// Let's create it

		// Get or create the member
		if member, ok := utils.SliceFind(dbMembers, func(m *model.Member) bool { return m.Equal(board.Member) }); ok {
			board.MemberID = member.ID
		} else {
			if err := c.memberRepo.Create(ctx, &board.Member); err != nil {
				return fmt.Errorf("creating member for new board %+v | %w", board, err)
			}
			dbMembers = append(dbMembers, &board.Member) // Update our db list to avoid creating duplicate members
			board.MemberID = board.Member.ID
		}

		// Get or create the year
		if year, ok := utils.SliceFind(dbYears, func(y *model.Year) bool { return y.Equal(board.Year) }); ok {
			board.YearID = year.ID
		} else {
			if err := c.yearRepo.Create(ctx, &board.Year); err != nil {
				return fmt.Errorf("creating new year for new board %+v | %w", board, err)
			}
			dbYears = append(dbYears, &board.Year) // Update our db list to avoid creating duplicate years
			board.YearID = board.Year.ID
		}

		if err := c.boardRepo.Create(ctx, &board); err != nil {
			return err
		}
	}

	// Refresh our database boards
	dbBoards, err = c.boardRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	// Delete old boards
	for _, board := range dbBoards {
		if !board.IsOrganizer {
			// Don't delete manually created board members
			// They are not on the website as board member
			// and are either in a development environment or an event admin
			continue
		}

		if ok := slices.ContainsFunc(websiteBoards, func(b model.Board) bool { return b.Equal(*board) }); !ok {
			if err := c.boardRepo.Delete(ctx, *board); err != nil {
				return err
			}
		}
	}

	return nil
}
