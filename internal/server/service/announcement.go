package service

import (
	"context"
	"time"

	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/server/dto"
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

func (a *Announcement) Save(ctx context.Context, announcementSave dto.Announcement) (dto.Announcement, error) {
	announcement := announcementSave.ToModel()

	if announcement.SendTime.Before(time.Now()) {
		return dto.Announcement{}, fiber.ErrBadRequest
	}

	event, err := a.event.GetByID(ctx, announcement.EventID)
	if err != nil {
		return dto.Announcement{}, fiber.ErrInternalServerError
	}
	if event == nil {
		return dto.Announcement{}, fiber.ErrBadRequest
	}

	if announcement.SendTime.After(event.StartTime) {
		return dto.Announcement{}, fiber.ErrBadRequest
	}

	if err = a.service.withRollback(ctx, func(ctx context.Context) error {
		update := false
		if announcement.ID == 0 {
			err = a.announcement.Create(ctx, &announcement)
		} else {
			err = a.announcement.Update(ctx, announcement)
			update = true
		}

		if err != nil {
			return err
		}

		if err = a.service.mattermost.ScheduleAnnouncement(ctx, announcement, update); err != nil {
			return err
		}

		return nil
	}); err != nil {
		zap.S().Error(err)
		return dto.Announcement{}, fiber.ErrInternalServerError
	}

	return dto.AnnouncementDTO(announcement), nil
}
