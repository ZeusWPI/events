package repository

import (
	"context"
	"fmt"

	"github.com/ZeusWPI/events/internal/pkg/db/sqlc"
	"github.com/ZeusWPI/events/internal/pkg/model"
	"github.com/ZeusWPI/events/pkg/db"
	"github.com/ZeusWPI/events/pkg/util"
)

// AcademicYear provides all model.AcademicYear related database operations
type AcademicYear interface {
	GetAll() ([]*model.AcademicYear, error)
	Save(*model.AcademicYear) error
}

type academicYearRepo struct {
	db db.DB
}

// Interface compliance
var _ AcademicYear = (*academicYearRepo)(nil)

// GetAll returns all academic year in desc order according to start year
func (r *academicYearRepo) GetAll() ([]*model.AcademicYear, error) {
	yearsDB, err := r.db.Queries().AcademicYearGetAll(context.Background())
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
func (r *academicYearRepo) Save(a *model.AcademicYear) error {
	var id int32
	var err error

	if a.ID == 0 {
		// Create
		id, err = r.db.Queries().AcademicYearCreate(context.Background(), sqlc.AcademicYearCreateParams{
			StartYear: int32(a.StartYear),
			EndYear:   int32(a.EndYear),
		})
	} else {
		// Update
		id = int32(a.ID)
		err = r.db.Queries().AcademicYearUpdate(context.Background(), sqlc.AcademicYearUpdateParams{
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
