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

func (a *Announcement) GetByYear(ctx context.Context, yearID int) ([]*model.Announcement, error) {
	announcements, err := a.repo.queries(ctx).AnnouncmentGetByYear(ctx, int32(yearID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get all announcements %w", err)
	}

	announcementMap := make(map[int32]*model.Announcement)

	// TODO: Cleanup
	for _, a := range announcements {
		if _, ok := announcementMap[a.ID]; !ok {
			err := ""
			if a.Error.Valid {
				err = a.Error.String
			}
			announcementMap[a.ID] = &model.Announcement{
				ID:       int(a.ID),
				EventIDs: []int{},
				YearID:   int(a.YearID),
				Content:  a.Content,
				SendTime: a.SendTime.Time,
				Send:     a.Send,
				Error:    err,
			}
		}

		if a.EventID.Valid {
			announcementMap[a.ID].EventIDs = append(announcementMap[a.ID].EventIDs, int(a.EventID.Int32))
		}
	}

	return utils.MapValues(announcementMap), nil
}

func (a *Announcement) GetByEvents(ctx context.Context, events []model.Event) ([]*model.Announcement, error) {
	announcements, err := a.repo.queries(ctx).AnnouncementGetByEvents(ctx, utils.SliceMap(events, func(e model.Event) int32 { return int32(e.ID) }))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get announcement by events %+v | %w", events, err)
	}

	announcementMap := make(map[int32]*model.Announcement)

	for _, a := range announcements {
		if _, ok := announcementMap[a.ID]; !ok {
			err := ""
			if a.Error.Valid {
				err = a.Error.String
			}
			announcementMap[a.ID] = &model.Announcement{
				ID:       int(a.ID),
				EventIDs: []int{},
				YearID:   int(a.YearID),
				Content:  a.Content,
				SendTime: a.SendTime.Time,
				Send:     a.Send,
				Error:    err,
			}
		}

		if a.EventID.Valid {
			announcementMap[a.ID].EventIDs = append(announcementMap[a.ID].EventIDs, int(a.EventID.Int32))
		}
	}

	return utils.MapValues(announcementMap), nil
}

func (a *Announcement) GetUnsend(ctx context.Context) ([]*model.Announcement, error) {
	announcements, err := a.repo.queries(ctx).AnnouncementGetUnsend(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get unsend announcements %w", err)
	}

	announcementMap := make(map[int32]*model.Announcement)

	for _, a := range announcements {
		if _, ok := announcementMap[a.ID]; !ok {
			err := ""
			if a.Error.Valid {
				err = a.Error.String
			}
			announcementMap[a.ID] = &model.Announcement{
				ID:       int(a.ID),
				EventIDs: []int{},
				YearID:   int(a.YearID),
				Content:  a.Content,
				SendTime: a.SendTime.Time,
				Send:     a.Send,
				Error:    err,
			}
		}

		if a.EventID.Valid {
			announcementMap[a.ID].EventIDs = append(announcementMap[a.ID].EventIDs, int(a.EventID.Int32))
		}
	}

	return utils.MapValues(announcementMap), nil
}

func (a *Announcement) Create(ctx context.Context, announcement *model.Announcement) error {
	return a.repo.WithRollback(ctx, func(ctx context.Context) error {
		id, err := a.repo.queries(ctx).AnnouncementCreate(ctx, sqlc.AnnouncementCreateParams{
			YearID:   int32(announcement.YearID),
			Content:  announcement.Content,
			SendTime: pgtype.Timestamptz{Valid: true, Time: announcement.SendTime},
			Send:     announcement.Send,
			Error:    pgtype.Text{Valid: announcement.Error != "", String: announcement.Error},
		})
		if err != nil {
			return fmt.Errorf("create announcement %+v | %w", *announcement, err)
		}

		announcement.ID = int(id)

		if len(announcement.EventIDs) > 0 {
			if err := a.repo.queries(ctx).AnnouncementEventCreateBatch(ctx, sqlc.AnnouncementEventCreateBatchParams{
				Column1: utils.SliceRepeat(id, len(announcement.EventIDs)),
				Column2: utils.SliceMap(announcement.EventIDs, func(id int) int32 { return int32(id) }),
			}); err != nil {
				return fmt.Errorf("create announcement events %+v | %w", *announcement, err)
			}
		}

		return nil
	})
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
