package repository

import (
	"context"
	"fmt"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/sqlc"
	"github.com/ZeusWPI/events/pkg/utils"
)

type Year struct {
	repo Repository
}

func newYear(repo Repository) *Year {
	return &Year{
		repo: repo,
	}
}

func (y *Year) GetAll(ctx context.Context) ([]*model.Year, error) {
	yearsDB, err := y.repo.queries(ctx).YearGetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all years | %w", err)
	}

	return utils.SliceMap(yearsDB, model.YearModel), nil
}

func (y *Year) GetLast(ctx context.Context) (*model.Year, error) {
	year, err := y.repo.queries(ctx).YearGetLast(ctx)
	if err != nil {
		return nil, fmt.Errorf("get last year %w", err)
	}

	return model.YearModel(year), nil
}

func (y *Year) Create(ctx context.Context, year *model.Year) error {
	id, err := y.repo.queries(ctx).YearCreate(ctx, sqlc.YearCreateParams{
		YearStart: int32(year.Start),
		YearEnd:   int32(year.End),
	})
	if err != nil {
		return fmt.Errorf("create year %+v | %w", *year, err)
	}

	year.ID = int(id)

	return nil
}
