package service

import (
	"context"

	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/server/dto"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Announcement struct {
	service Service

	announcement repository.Announcement
}

func (s *Service) NewAnnouncement() *Announcement {
	return &Announcement{
		service:      *s,
		announcement: *s.repo.NewAnnouncement(),
	}
}

func (a *Announcement) Save(ctx context.Context, announcementSave dto.Announcement) (dto.Announcement, error) {
	announcement := announcementSave.ToModel()

	var err error

	if announcement.ID == 0 {
		err = a.announcement.Create(ctx, &announcement)
	} else {
		err = a.announcement.Update(ctx, announcement)
	}

	if err != nil {
		zap.S().Error(err)
		return dto.Announcement{}, fiber.ErrInternalServerError
	}

	return dto.AnnouncementDTO(announcement), nil
}
