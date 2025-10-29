package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/sqlc"
	"github.com/ZeusWPI/events/pkg/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

type Board struct {
	repo Repository

	year   Year
	member Member
}

func (r *Repository) NewBoard() *Board {
	return &Board{
		repo:   *r,
		year:   *r.NewYear(),
		member: *r.NewMember(),
	}
}

func (b *Board) GetByMemberYear(ctx context.Context, memberID int, yearID int) (*model.Board, error) {
	board, err := b.repo.queries(ctx).BoardGetByMemberYear(ctx, sqlc.BoardGetByMemberYearParams{
		ID:   int32(memberID),
		ID_2: int32(yearID),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get board by member and year %d | %d | %w", memberID, yearID, err)
	}

	return model.BoardModel(board.Board, board.Member, board.Year), nil
}

func (b *Board) GetAll(ctx context.Context) ([]*model.Board, error) {
	boards, err := b.repo.queries(ctx).BoardGetAll(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get all boards %w", err)
	}

	return utils.SliceMap(boards, func(b sqlc.BoardGetAllRow) *model.Board {
		return model.BoardModel(b.Board, b.Member, b.Year)
	}), nil
}

func (b *Board) GetByIDs(ctx context.Context, boardIDs []int) ([]*model.Board, error) {
	boards, err := b.repo.queries(ctx).BoardGetByIds(ctx, utils.SliceMap(boardIDs, func(id int) int32 { return int32(id) }))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get board by ids %+v | %w", boardIDs, err)
	}

	return utils.SliceMap(boards, func(b sqlc.BoardGetByIdsRow) *model.Board {
		return model.BoardModel(b.Board, b.Member, b.Year)
	}), nil
}

func (b *Board) GetByYear(ctx context.Context, yearID int) ([]*model.Board, error) {
	boards, err := b.repo.queries(ctx).BoardGetByYear(ctx, int32(yearID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get all populated boards by year %w", err)
	}

	return utils.SliceMap(boards, func(b sqlc.BoardGetByYearRow) *model.Board {
		return model.BoardModel(b.Board, b.Member, b.Year)
	}), nil
}

func (b *Board) Create(ctx context.Context, board *model.Board) error {
	id, err := b.repo.queries(ctx).BoardCreate(ctx, sqlc.BoardCreateParams{
		MemberID:    int32(board.MemberID),
		YearID:      int32(board.YearID),
		Role:        board.Role,
		IsOrganizer: board.IsOrganizer,
		Mattermost:  pgtype.Text{String: board.Mattermost, Valid: board.Mattermost != ""},
	})
	if err != nil {
		return fmt.Errorf("create board %+v | %w", *board, err)
	}

	board.ID = int(id)

	return nil
}

func (b *Board) Update(ctx context.Context, board model.Board) error {
	if err := b.repo.queries(ctx).BoardUpdate(ctx, sqlc.BoardUpdateParams{
		ID:          int32(board.ID),
		MemberID:    int32(board.MemberID),
		YearID:      int32(board.YearID),
		Role:        board.Role,
		IsOrganizer: board.IsOrganizer,
		Mattermost:  pgtype.Text{String: board.Mattermost, Valid: board.Mattermost != ""},
	}); err != nil {
		return fmt.Errorf("update board %+v | %w", board, err)
	}

	return nil
}

func (b *Board) Delete(ctx context.Context, board model.Board) error {
	return b.repo.WithRollback(ctx, func(ctx context.Context) error {
		if err := b.repo.queries(ctx).BoardDelete(ctx, int32(board.ID)); err != nil {
			return fmt.Errorf("Delete board %+v | %w", board, err)
		}

		boards, err := b.repo.queries(ctx).BoardGetByMember(ctx, int32(board.MemberID))
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("get board by member id %+v | %w", board, err)
		}

		if len(boards) == 0 || (err != nil && errors.Is(err, sql.ErrNoRows)) {
			// No more board entries, also delete the member
			if err := b.member.Delete(ctx, board.MemberID); err != nil {
				return err
			}
		}

		return nil
	})
}
