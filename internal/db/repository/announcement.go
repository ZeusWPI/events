package repository

import (
	"context"
	"fmt"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

// Announcement provides all model.Announcement related database operations
type Announcement interface {
	GetAll(context.Context) ([]*model.Announcement, error)
	GetAllByEvent(context.Context, model.Event) ([]*model.Announcement, error)
	Save(context.Context, *model.Announcement) error
	Delete(context.Context, *model.Announcement) error
}

type announcementRepo struct {
	repo Repository
}

// Interface compliance
var _ Announcement = (*announcementRepo)(nil)

// GetAll returns all announcements
func (r *announcementRepo) GetAll(ctx context.Context) ([]*model.Announcement, error) {
	announcements, err := r.repo.queries(ctx).AnnouncementGetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get all events %v", err)
	}

	announcementModels := make([]*model.Announcement, 0, len(announcements))

	for _, a := range announcements {
		target, err := model.FromTargetString(a.Target.String)
		if err != nil {
			return nil, err
		}

		announcementModels = append(announcementModels, &model.Announcement{
			ID:      int(a.ID),
			Content: a.Content,
			Time:    a.Time.Time,
			Target:  target,
			Member: model.Member{
				ID:       int(a.ID_2),
				ZauthID:  int(a.ZauthID.Int32),
				Name:     a.Name,
				Username: a.Username.String,
			},
		})
	}

	return announcementModels, nil
}

// GetAllByEvent returns all announcements associated with an event
func (r *announcementRepo) GetAllByEvent(ctx context.Context, event model.Event) ([]*model.Announcement, error) {
	announcements, err := r.repo.queries(ctx).AnnouncementGetByEvent(ctx, int32(event.ID))
	if err != nil {
		return nil, fmt.Errorf("unable to get all announcements by event %+v | %v", event, err)
	}

	announcementModels := make([]*model.Announcement, 0, len(announcements))

	for _, a := range announcements {
		target, err := model.FromTargetString(a.Target.String)
		if err != nil {
			return nil, err
		}

		announcementModels = append(announcementModels, &model.Announcement{
			ID:      int(a.ID),
			Content: a.Content,
			Time:    a.Time.Time,
			Target:  target,
			Member: model.Member{
				ID:       int(a.ID_2),
				ZauthID:  int(a.ZauthID.Int32),
				Name:     a.Name,
				Username: a.Username.String,
			},
		})
	}

	return announcementModels, nil
}

// Save creates a new announcement or updates an existing one
func (r *announcementRepo) Save(ctx context.Context, a *model.Announcement) error {
	var id int32
	var err error

	if a.ID == 0 {
		// Create
		id, err = r.repo.queries(ctx).AnnouncementCreate(ctx, sqlc.AnnouncementCreateParams{
			Content: a.Content,
			Time:    pgtype.Timestamptz{Time: a.Time, Valid: true},
			Target:  pgtype.Text{String: a.Target.String(), Valid: true},
			Event:   int32(a.Event.ID),
			Member:  pgtype.Int4{Int32: int32(a.Member.ID), Valid: a.Member.ID != 0},
		})
	} else {
		// Updates
		id = int32(a.ID)
		err = r.repo.queries(ctx).AnnouncementUpdate(ctx, sqlc.AnnouncementUpdateParams{
			ID:      int32(a.ID),
			Content: a.Content,
			Time:    pgtype.Timestamptz{Time: a.Time, Valid: true},
			Target:  pgtype.Text{String: a.Target.String(), Valid: true},
			Event:   int32(a.Event.ID),
			Member:  pgtype.Int4{Int32: int32(a.Member.ID), Valid: a.Member.ID != 0},
		})
	}

	if err != nil {
		return fmt.Errorf("unable to save announcement %+v | %v", *a, err)
	}

	a.ID = int(id)

	return nil
}

// Delete deletes an announcement
func (r *announcementRepo) Delete(ctx context.Context, e *model.Announcement) error {
	if e.ID == 0 {
		return fmt.Errorf("Announcement has no ID %+v", *e)
	}

	if err := r.repo.queries(ctx).AnnouncementDelete(ctx, int32(e.ID)); err != nil {
		return fmt.Errorf("unable to delete announcement %+v | %v", *e, err)
	}

	return nil
}
