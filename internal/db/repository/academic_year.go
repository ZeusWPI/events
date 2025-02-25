package repository

import (
	"context"
	"fmt"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/sqlc"
	"github.com/ZeusWPI/events/pkg/util"
)

// AcademicYear provides all model.AcademicYear related database operations
type AcademicYear interface {
	GetAll(context.Context) ([]*model.AcademicYear, error)
	Save(context.Context, *model.AcademicYear) error
}

type academicYearRepo struct {
	repo Repository
}

// Interface compliance
var _ AcademicYear = (*academicYearRepo)(nil)

// GetAll returns all academic year in desc order according to start year
func (r *academicYearRepo) GetAll(ctx context.Context) ([]*model.AcademicYear, error) {
	yearsDB, err := r.repo.queries(ctx).AcademicYearGetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("Unable to get all academic years | %v", err)
	}

	return util.SliceMap(yearsDB, func(y sqlc.AcademicYear) *model.AcademicYear {
		return &model.AcademicYear{
			ID:        int(y.ID),
			StartYear: int(y.StartYear),
			EndYear:   int(y.EndYear),
		}
	}), nil
}

// Save creates a new academic year or updates an existing one
func (r *academicYearRepo) Save(ctx context.Context, a *model.AcademicYear) error {
	var id int32
	var err error

	if a.ID == 0 {
		// Create
		id, err = r.repo.queries(ctx).AcademicYearCreate(ctx, sqlc.AcademicYearCreateParams{
			StartYear: int32(a.StartYear),
			EndYear:   int32(a.EndYear),
		})
	} else {
		// Update
		id = int32(a.ID)
		err = r.repo.queries(ctx).AcademicYearUpdate(ctx, sqlc.AcademicYearUpdateParams{
			ID:        int32(a.ID),
			StartYear: int32(a.StartYear),
			EndYear:   int32(a.EndYear),
		})
	}

	if err != nil {
		return fmt.Errorf("Unable to save academic year %+v | %v", *a, err)
	}

	a.ID = int(id)

	return nil
}
