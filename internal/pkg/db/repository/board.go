package repository

import (
	"context"
	"fmt"

	"github.com/ZeusWPI/events/internal/pkg/db/sqlc"
	"github.com/ZeusWPI/events/internal/pkg/model"
	"github.com/ZeusWPI/events/pkg/db"
	"github.com/ZeusWPI/events/pkg/util"
)

// Board provides all model.Board related database operations
type Board interface {
	GetAll() ([]*model.Board, error)
	Save(*model.Board) error
}

type boardRepo struct {
	db db.DB

	year   AcademicYear
	member Member
}

// Interface compliance
var _ Board = (*boardRepo)(nil)

// GetAll returns all boards
func (r *boardRepo) GetAll() ([]*model.Board, error) {
	boards, err := r.db.Queries().BoardGetAll(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Unable to get all boards | %v", err)
	}

	return util.SliceMap(boards, func(b sqlc.BoardGetAllRow) *model.Board {
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
			AcademicYear: model.AcademicYear{
				ID:        int(b.ID_3),
				StartYear: int(b.StartYear),
				EndYear:   int(b.EndYear),
			},
			Role: b.Role,
		}
	}), nil
}

// Save creates a new board
func (r *boardRepo) Save(b *model.Board) error {
	if b.ID != 0 {
		// Already in database
		return nil
	}

	return r.db.WithRollback(context.Background(), func(q *sqlc.Queries) error {
		if b.Member.ID == 0 {
			err := r.member.Save(&b.Member)
			if err != nil {
				return err
			}
		}

		if b.AcademicYear.ID == 0 {
			err := r.year.Save(&b.AcademicYear)
			if err != nil {
				return err
			}
		}

		id, err := q.BoardCreate(context.Background(), sqlc.BoardCreateParams{
			Member:       int32(b.Member.ID),
			AcademicYear: int32(b.AcademicYear.ID),
			Role:         b.Role,
		})
		if err != nil {
			return fmt.Errorf("Unable to save board %+v | %v", *b, err)
		}

		b.ID = int(id)

		return nil
	})
}
