package repository

import (
	"context"
	"fmt"

	"github.com/ZeusWPI/events/internal/db/sqlc"
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

func (o *Organizer) DeleteByBoardEvent(ctx context.Context, boardID, eventID int) error {
	if err := o.repo.queries(ctx).OrganizerDeleteByBoardEvent(ctx, sqlc.OrganizerDeleteByBoardEventParams{BoardID: int32(boardID), EventID: int32(eventID)}); err != nil {
		return fmt.Errorf("delete organizer %w", err)
	}

	return nil
}
