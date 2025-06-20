package service

import (
	"context"
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/server/dto"
	"github.com/ZeusWPI/events/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Mail struct {
	service Service

	mail repository.Mail
}

func (s *Service) NewMail() *Mail {
	return &Mail{
		service: *s,
		mail:    *s.repo.NewMail(),
	}
}

func (m *Mail) GetAll(ctx context.Context) ([]dto.Mail, error) {
	mails, err := m.mail.GetAll(ctx)
	if err != nil {
		zap.S().Error(err)
		return []dto.Mail{}, fiber.ErrInternalServerError
	}

	return utils.SliceMap(mails, dto.MailDTO), nil
}

func (m *Mail) Save(ctx context.Context, mailSave dto.MailSave) (dto.Mail, error) {
	mail := model.Mail{
		ID:       mailSave.ID,
		Content:  mailSave.Content,
		SendTime: mailSave.SendTime,
		Send:     false,
		Error:    "",
	}

	if mail.SendTime.Before(time.Now()) {
		return dto.Mail{}, fiber.ErrBadRequest
	}

	var err error
	update := false
	if mail.ID == 0 {
		err = m.mail.Create(ctx, &mail, mailSave.EventIDs)
	} else {
		err = m.mail.Update(ctx, mail, mailSave.EventIDs)
		update = true
	}

	if err != nil {
		zap.S().Error(err)
		return dto.Mail{}, fiber.ErrInternalServerError
	}

	if err = m.service.mail.ScheduleMailAll(ctx, mail, update); err != nil {
		zap.S().Error(err)
		return dto.Mail{}, fiber.ErrInternalServerError
	}

	return dto.MailDTO(&mail), nil
}
