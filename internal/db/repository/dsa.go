package repository

import (
	"context"
	"fmt"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/sqlc"
	"github.com/ZeusWPI/events/pkg/utils"
)

type DSA struct {
	repo Repository
}

func (r *Repository) NewDSA() *DSA {
	return &DSA{
		repo: *r,
	}
}

func (d *DSA) GetByEvents(ctx context.Context, events []model.Event) ([]*model.DSA, error) {
	dsa, err := d.repo.queries(ctx).DsaGetByEvents(ctx, utils.SliceMap(events, func(e model.Event) int32 { return int32(e.ID) }))
	if err != nil {
		return nil, fmt.Errorf("get all dsa by events %+v | %w", events, err)
	}

	return utils.SliceMap(dsa, model.DSAModel), nil
}

func (d *DSA) Create(ctx context.Context, dsa *model.DSA) error {
	id, err := d.repo.queries(ctx).DsaCreate(ctx, sqlc.DsaCreateParams{
		EventID: int32(dsa.EventID),
		Entry:   dsa.Entry,
	})
	if err != nil {
		return fmt.Errorf("create dsa %+v | %w", *dsa, err)
	}

	dsa.ID = int(id)

	return nil
}

func (d *DSA) Delete(ctx context.Context, dsaID int) error {
	if err := d.repo.queries(ctx).DsaDelete(ctx, int32(dsaID)); err != nil {
		return fmt.Errorf("delete dsa %d | %w", dsaID, err)
	}

	return nil
}
