package repository

import (
	"context"
	"fmt"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/sqlc"
	"github.com/ZeusWPI/events/pkg/util"
)

// Board provides all model.Board related database operations
type Board interface {
	GetAllWithMemberYear(context.Context) ([]*model.Board, error)
	Save(context.Context, *model.Board) error
}

type boardRepo struct {
	repo Repository

	year   Year
	member Member
}

// Interface compliance
var _ Board = (*boardRepo)(nil)

// GetAll returns all boards
func (r *boardRepo) GetAllWithMemberYear(ctx context.Context) ([]*model.Board, error) {
	boards, err := r.repo.queries(ctx).BoardGetAllWithMemberYear(ctx)
	if err != nil {
		return nil, fmt.Errorf("Unable to get all boards | %v", err)
	}

	return util.SliceMap(boards, func(b sqlc.BoardGetAllWithMemberYearRow) *model.Board {
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
				ID:        int(b.ID_3),
				StartYear: int(b.StartYear),
				EndYear:   int(b.EndYear),
			},
			Role: b.Role,
		}
	}), nil
}

// Save creates a new board
func (r *boardRepo) Save(ctx context.Context, b *model.Board) error {
	if b.ID != 0 {
		// Already in database
		return nil
	}

	return r.repo.withRollback(ctx, func(c context.Context) error {
		if b.Member.ID == 0 {
			err := r.member.Save(c, &b.Member)
			if err != nil {
				return err
			}
		}

		if b.Year.ID == 0 {
			err := r.year.Save(c, &b.Year)
			if err != nil {
				return err
			}
		}

		id, err := r.repo.queries(c).BoardCreate(c, sqlc.BoardCreateParams{
			Member: int32(b.Member.ID),
			Year:   int32(b.Year.ID),
			Role:   b.Role,
		})
		if err != nil {
			return fmt.Errorf("Unable to save board %+v | %v", *b, err)
		}

		b.ID = int(id)

		return nil
	})
}
