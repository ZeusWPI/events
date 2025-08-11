package service

import (
	"context"
	"time"

	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/server/dto"
	"github.com/ZeusWPI/events/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Announcement struct {
	service Service

	announcement repository.Announcement
	event        repository.Event
}

func (s *Service) NewAnnouncement() *Announcement {
	return &Announcement{
		service:      *s,
		announcement: *s.repo.NewAnnouncement(),
		event:        *s.repo.NewEvent(),
	}
}

func (a *Announcement) GetByYear(ctx context.Context, yearID int) ([]dto.Announcement, error) {
	announcements, err := a.announcement.GetByYear(ctx, yearID)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}

	return utils.SliceMap(announcements, dto.AnnouncementDTO), nil
}

func (a *Announcement) Save(ctx context.Context, announcementSave dto.Announcement) (dto.Announcement, error) {
	announcement := announcementSave.ToModel()

	if announcement.SendTime.Before(time.Now()) {
		return dto.Announcement{}, fiber.ErrBadRequest
	}

	events, err := a.event.GetByIDs(ctx, announcement.EventIDs)
	if err != nil {
		return dto.Announcement{}, fiber.ErrInternalServerError
	}
	if len(events) != len(announcement.EventIDs) {
		return dto.Announcement{}, fiber.ErrBadRequest
	}

	for _, event := range events {
		if announcement.SendTime.After(event.StartTime) {
			return dto.Announcement{}, fiber.ErrBadRequest
		}
	}

	if announcement.ID != 0 {
		oldAnnouncement, err := a.announcement.GetByID(ctx, announcement.ID)
		if err != nil {
			zap.S().Error(err)
			return dto.Announcement{}, fiber.ErrInternalServerError
		}
		if oldAnnouncement == nil {
			return dto.Announcement{}, fiber.ErrBadRequest
		}

		if oldAnnouncement.Send || oldAnnouncement.Error != "" {
			return dto.Announcement{}, fiber.ErrBadRequest
		}
	}

	if err = a.service.withRollback(ctx, func(ctx context.Context) error {
		update := false
		if announcement.ID == 0 {
			err = a.announcement.Create(ctx, announcement)
		} else {
			err = a.announcement.Update(ctx, *announcement)
			update = true
		}

		if err != nil {
			return err
		}

		if err = a.service.mattermost.ScheduleAnnouncement(ctx, *announcement, update); err != nil {
			return err
		}

		return nil
	}); err != nil {
		zap.S().Error(err)
		return dto.Announcement{}, fiber.ErrInternalServerError
	}

	return dto.AnnouncementDTO(announcement), nil
}

func (a *Announcement) Delete(ctx context.Context, announcementID int) error {
	announcement, err := a.announcement.GetByID(ctx, announcementID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}
	if announcement == nil {
		return fiber.ErrBadRequest
	}

	if announcement.Send || announcement.Error != "" {
		return fiber.ErrBadRequest
	}

	if err := a.announcement.Delete(ctx, announcementID); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return nil
}
