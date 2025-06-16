package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/sqlc"
	"github.com/ZeusWPI/events/pkg/utils"
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

func (b *Board) GetAll(ctx context.Context) ([]*model.Board, error) {
	boards, err := b.repo.queries(ctx).BoardGetAll(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get all boards %w", err)
	}

	return utils.SliceMap(boards, model.BoardModel), nil
}

func (b *Board) GetByYearPopulated(ctx context.Context, yearID int) ([]*model.Board, error) {
	boards, err := b.repo.queries(ctx).BoardGetByYearPopulated(ctx, int32(yearID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get all populated boards by year %w", err)
	}

	return utils.SliceMap(boards, func(b sqlc.BoardGetByYearPopulatedRow) *model.Board {
		username := ""
		if b.Username.Valid {
			username = b.Username.String
		}

		return &model.Board{
			ID: int(b.ID),
			Member: model.Member{
				ID:       int(b.ID_2),
				Name:     b.Name,
				Username: username,
			},
			Year: model.Year{
				ID:    int(b.ID_3),
				Start: int(b.YearStart),
				End:   int(b.YearEnd),
			},
			Role: b.Role,
		}
	}), nil
}

func (b *Board) GetByMemberYear(ctx context.Context, member model.Member, year model.Year) (*model.Board, error) {
	board, err := b.repo.queries(ctx).BoardGetByMemberYear(ctx, sqlc.BoardGetByMemberYearParams{
		ID:   int32(member.ID),
		ID_2: int32(year.ID),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get board by member and year %+v | %+v | %w", member, year, err)
	}

	username := ""
	if board.Username.Valid {
		username = board.Username.String
	}

	return &model.Board{
		ID: int(board.ID),
		Member: model.Member{
			ID:       int(board.ID_2),
			Name:     board.Name,
			Username: username,
		},
		Year: model.Year{
			ID:    int(board.ID_3),
			Start: int(board.YearStart),
			End:   int(board.YearEnd),
		},
		Role: board.Role,
	}, nil
}

func (b *Board) Create(ctx context.Context, board *model.Board) error {
	id, err := b.repo.queries(ctx).BoardCreate(ctx, sqlc.BoardCreateParams{
		MemberID: int32(board.MemberID),
		YearID:   int32(board.YearID),
		Role:     board.Role,
	})
	if err != nil {
		return fmt.Errorf("create board %+v | %w", *board, err)
	}

	board.ID = int(id)

	return nil
}

func (b *Board) Delete(ctx context.Context, board model.Board) error {
	return b.repo.WithRollback(ctx, func(ctx context.Context) error {
		if err := b.repo.queries(ctx).BoardDelete(ctx, int32(board.ID)); err != nil {
			return fmt.Errorf("Delete board %+v | %w", board, err)
		}

		boards, err := b.repo.queries(ctx).BoardGetByMemberID(ctx, int32(board.MemberID))
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
