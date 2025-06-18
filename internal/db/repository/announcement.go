package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/sqlc"
	"github.com/ZeusWPI/events/pkg/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

type Announcement struct {
	repo Repository
}

func (r *Repository) NewAnnouncement() *Announcement {
	return &Announcement{
		repo: *r,
	}
}

func (a *Announcement) GetByEvents(ctx context.Context, events []model.Event) ([]*model.Announcement, error) {
	announcements, err := a.repo.queries(ctx).AnnouncementGetByEvents(ctx, utils.SliceMap(events, func(e model.Event) int32 { return int32(e.ID) }))
	if err != nil {
		return nil, fmt.Errorf("get announcement by events %+v | %w", events, err)
	}

	return utils.SliceMap(announcements, model.AnnouncementModel), nil
}

func (a *Announcement) GetUnsend(ctx context.Context) ([]*model.Announcement, error) {
	announcements, err := a.repo.queries(ctx).AnnouncementGetUnsend(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get unsend announcements %w", err)
	}

	return utils.SliceMap(announcements, model.AnnouncementModel), nil
}

func (a *Announcement) Create(ctx context.Context, announcement *model.Announcement) error {
	id, err := a.repo.queries(ctx).AnnouncementCreate(ctx, sqlc.AnnouncementCreateParams{
		EventID:  int32(announcement.EventID),
		Content:  announcement.Content,
		SendTime: pgtype.Timestamptz{Valid: true, Time: announcement.SendTime},
		Send:     announcement.Send,
	})
	if err != nil {
		return fmt.Errorf("create announcement %+v | %w", *announcement, err)
	}

	announcement.ID = int(id)

	return nil
}

func (a *Announcement) Update(ctx context.Context, announcement model.Announcement) error {
	if err := a.repo.queries(ctx).AnnouncementUpdate(ctx, sqlc.AnnouncementUpdateParams{
		ID:       int32(announcement.ID),
		Content:  announcement.Content,
		SendTime: pgtype.Timestamptz{Valid: true, Time: announcement.SendTime},
	}); err != nil {
		return fmt.Errorf("update announcement %+v | %w", announcement, err)
	}

	return nil
}

func (a *Announcement) Send(ctx context.Context, announcementID int) error {
	if err := a.repo.queries(ctx).AnnouncementSend(ctx, int32(announcementID)); err != nil {
		return fmt.Errorf("send announcement %d | %w", announcementID, err)
	}

	return nil
}

func (a *Announcement) Error(ctx context.Context, announcement model.Announcement) error {
	if err := a.repo.queries(ctx).AnnouncementError(ctx, sqlc.AnnouncementErrorParams{
		ID:    int32(announcement.ID),
		Error: pgtype.Text{Valid: true, String: announcement.Error},
	}); err != nil {
		return fmt.Errorf("error announcement %+v | %w", announcement, err)
	}

	return nil
}
