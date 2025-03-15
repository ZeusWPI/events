package repository

import (
	"context"
	"fmt"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/sqlc"
	"github.com/ZeusWPI/events/pkg/util"
)

// Organizer provides all model.Organizer related database operations
type Organizer interface {
	GetByYearWithBoard(context.Context, model.Year) ([]*model.Organizer, error)
	Save(context.Context, *model.Organizer) error
	Delete(context.Context, model.Organizer) error
}

type organizerRepo struct {
	repo Repository
}

var _ Organizer = (*organizerRepo)(nil)

// GetByYearWithBoard returns all organizers of a given year with the board field completely populated
func (r *organizerRepo) GetByYearWithBoard(ctx context.Context, year model.Year) ([]*model.Organizer, error) {
	organizers, err := r.repo.queries(ctx).OrganizerGetByYearWithBoard(ctx, int32(year.ID))
	if err != nil {
		return nil, fmt.Errorf("unable to get all organizers by year %+v | %v", year, err)
	}

	return util.SliceMap(organizers, func(o sqlc.OrganizerGetByYearWithBoardRow) *model.Organizer {
		return &model.Organizer{
			ID: int(o.ID),
			Event: model.Event{
				ID: int(o.Event),
			},
			Board: model.Board{
				ID:   int(o.Board),
				Role: o.Role,
				Member: model.Member{
					ID:       int(o.Member),
					Name:     o.Name,
					Username: o.Username.String,
				},
				Year: model.Year{
					ID: int(o.Year),
				},
			},
		}
	}), nil
}

func (r *organizerRepo) Save(ctx context.Context, organizer *model.Organizer) error {
	id, err := r.repo.queries(ctx).OrganizerCreate(ctx, sqlc.OrganizerCreateParams{Event: int32(organizer.Event.ID), Board: int32(organizer.Board.ID)})
	if err != nil {
		return fmt.Errorf("unable to save organizer %+v | %v", *organizer, err)
	}

	organizer.ID = int(id)

	return nil
}

func (r *organizerRepo) Delete(ctx context.Context, organizer model.Organizer) error {
	if err := r.repo.queries(ctx).OrganizerDelete(ctx, int32(organizer.ID)); err != nil {
		return fmt.Errorf("unable to delete organizer %+v | %v", organizer, err)
	}

	return nil
}
