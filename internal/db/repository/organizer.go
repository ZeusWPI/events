package repository

import (
	"context"
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

func (o *Organizer) Create(ctx context.Context, boardID, eventID int) error {
	if _, err := o.repo.queries(ctx).OrganizerCreate(ctx, sqlc.OrganizerCreateParams{
		BoardID: int32(boardID),
		EventID: int32(eventID),
	}); err != nil {
		return fmt.Errorf("create organizer %w", err)
	}

	return nil
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

func (o *Organizer) DeleteByBoardEvent(ctx context.Context, boardID, eventID int) error {
	if err := o.repo.queries(ctx).OrganizerDeleteByBoardEvent(ctx, sqlc.OrganizerDeleteByBoardEventParams{BoardID: int32(boardID), EventID: int32(eventID)}); err != nil {
		return fmt.Errorf("delete organizer %w", err)
	}

	return nil
}

func (o *Organizer) DeleteByEvent(ctx context.Context, eventID int) error {
	if err := o.repo.queries(ctx).OrganizerDeleteByEvent(ctx, int32(eventID)); err != nil {
		return fmt.Errorf("delete organizer by event %d | %w", eventID, err)
	}

	return nil
}
