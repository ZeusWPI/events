package repository

import (
	"context"
	"fmt"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/sqlc"
	"github.com/ZeusWPI/events/pkg/util"
)

// Year provides all model.Year related database operations
type Year interface {
	GetAll(context.Context) ([]*model.Year, error)
	GetLatest(context.Context) (*model.Year, error)
	Save(context.Context, *model.Year) error
}

type yearRepo struct {
	repo Repository
}

// Interface compliance
var _ Year = (*yearRepo)(nil)

// GetAll returns all year in desc order according to start year
func (r *yearRepo) GetAll(ctx context.Context) ([]*model.Year, error) {
	yearsDB, err := r.repo.queries(ctx).YearGetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get all years | %w", err)
	}

	return util.SliceMap(yearsDB, func(y sqlc.Year) *model.Year {
		return &model.Year{
			ID:        int(y.ID),
			StartYear: int(y.StartYear),
			EndYear:   int(y.EndYear),
		}
	}), nil
}

func (r *yearRepo) GetLatest(ctx context.Context) (*model.Year, error) {
	year, err := r.repo.queries(ctx).YearGetLatest(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get latest year %w", err)
	}

	return &model.Year{
		ID:        int(year.ID),
		StartYear: int(year.StartYear),
		EndYear:   int(year.EndYear),
	}, nil
}

// Save creates a new year or updates an existing one
func (r *yearRepo) Save(ctx context.Context, a *model.Year) error {
	var id int32
	var err error

	if a.ID == 0 {
		// Create
		id, err = r.repo.queries(ctx).YearCreate(ctx, sqlc.YearCreateParams{
			StartYear: int32(a.StartYear),
			EndYear:   int32(a.EndYear),
		})
	} else {
		// Update
		id = int32(a.ID)
		err = r.repo.queries(ctx).YearUpdate(ctx, sqlc.YearUpdateParams{
			ID:        int32(a.ID),
			StartYear: int32(a.StartYear),
			EndYear:   int32(a.EndYear),
		})
	}

	if err != nil {
		return fmt.Errorf("unable to save year %+v | %w", *a, err)
	}

	a.ID = int(id)

	return nil
}
