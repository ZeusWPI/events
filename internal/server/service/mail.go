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

type Mail struct {
	service Service

	events repository.Event
	mail   repository.Mail
}

func (s *Service) NewMail() *Mail {
	return &Mail{
		service: *s,
		events:  *s.repo.NewEvent(),
		mail:    *s.repo.NewMail(),
	}
}

func (m *Mail) GetByYear(ctx context.Context, yearID int) ([]dto.Mail, error) {
	mails, err := m.mail.GetByYear(ctx, yearID)
	if err != nil {
		zap.S().Error(err)
		return []dto.Mail{}, fiber.ErrInternalServerError
	}

	return utils.SliceMap(mails, dto.MailDTO), nil
}

func (m *Mail) Save(ctx context.Context, mailSave dto.Mail) (dto.Mail, error) {
	mail := mailSave.ToModel()

	if mail.SendTime.Before(time.Now()) {
		return dto.Mail{}, fiber.ErrBadRequest
	}

	events, err := m.events.GetByIDs(ctx, mail.EventIDs)
	if err != nil {
		return dto.Mail{}, fiber.ErrInternalServerError
	}
	if len(events) != len(mail.EventIDs) {
		return dto.Mail{}, fiber.ErrBadRequest
	}

	for _, event := range events {
		if mail.SendTime.After(event.StartTime) {
			return dto.Mail{}, fiber.ErrBadRequest
		}
	}

	if mail.ID != 0 {
		oldMail, err := m.mail.GetByID(ctx, mail.ID)
		if err != nil {
			zap.S().Error(err)
			return dto.Mail{}, fiber.ErrInternalServerError
		}
		if oldMail == nil {
			return dto.Mail{}, fiber.ErrBadRequest
		}

		if oldMail.Send || oldMail.Error != "" {
			return dto.Mail{}, fiber.ErrBadRequest
		}
	}

	if err = m.service.withRollback(ctx, func(ctx context.Context) error {
		update := false
		if mail.ID == 0 {
			err = m.mail.Create(ctx, mail)
		} else {
			err = m.mail.Update(ctx, *mail)
			update = true
		}

		if err != nil {
			return err
		}

		if err = m.service.mail.ScheduleMailAll(ctx, *mail, update); err != nil {
			return err
		}

		return nil
	}); err != nil {
		zap.S().Error(err)
		return dto.Mail{}, fiber.ErrInternalServerError
	}

	return dto.MailDTO(mail), nil
}

func (m *Mail) Delete(ctx context.Context, mailID int) error {
	mail, err := m.mail.GetByID(ctx, mailID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}
	if mail == nil {
		return fiber.ErrBadRequest
	}

	if mail.Send || mail.Error != "" {
		return fiber.ErrBadRequest
	}

	if err := m.mail.Delete(ctx, mailID); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return nil
}
