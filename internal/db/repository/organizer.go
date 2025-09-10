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

type Organizer struct {
	repo Repository
}

func (r *Repository) NewOrganizer() *Organizer {
	return &Organizer{
		repo: *r,
	}
}

func (o *Organizer) GetByEvents(ctx context.Context, events []model.Event) ([]*model.Organizer, error) {
	organizers, err := o.repo.queries(ctx).OrganizerGetByEvents(ctx, utils.SliceMap(events, func(event model.Event) int32 { return int32(event.ID) }))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get organizers by events %+v | %w", events, err)
	}

	return utils.SliceMap(organizers, func(o sqlc.OrganizerGetByEventsRow) *model.Organizer {
		return model.OrganizerModel(o.Organizer, o.Board, o.Event, o.Member, o.Year)
	}), nil
}

func (o *Organizer) CreateBatch(ctx context.Context, organizers []model.Organizer) error {
	if err := o.repo.queries(ctx).OrganizerCreateBatch(ctx, sqlc.OrganizerCreateBatchParams{
		Column1: utils.SliceMap(organizers, func(o model.Organizer) int32 { return int32(o.EventID) }),
		Column2: utils.SliceMap(organizers, func(o model.Organizer) int32 { return int32(o.BoardID) }),
	}); err != nil {
		return fmt.Errorf("create organizers batch %+v | %w", organizers, err)
	}

	return nil
}

func (o *Organizer) DeleteByEvent(ctx context.Context, eventID int) error {
	if err := o.repo.queries(ctx).OrganizerDeleteByEvent(ctx, int32(eventID)); err != nil {
		return fmt.Errorf("delete organizer by event %d | %w", eventID, err)
	}

	return nil
}
