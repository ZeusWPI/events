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

type Mail struct {
	repo Repository

	mailEvent MailEvent
}

func (r *Repository) NewMail() *Mail {
	return &Mail{
		repo:      *r,
		mailEvent: *r.NewMailEvent(),
	}
}

func (m *Mail) GetAll(ctx context.Context) ([]*model.Mail, error) {
	mails, err := m.repo.queries(ctx).MailGetAll(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get all mails %w", err)
	}

	return utils.SliceMap(mails, model.MailModel), nil
}

func (m *Mail) GetUnsend(ctx context.Context) ([]*model.Mail, error) {
	mails, err := m.repo.queries(ctx).MailGetUnsend(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get unsend mails %w", err)
	}

	return utils.SliceMap(mails, model.MailModel), nil
}

func (m *Mail) Create(ctx context.Context, mail *model.Mail, eventIDs []int) error {
	return m.repo.WithRollback(ctx, func(ctx context.Context) error {
		id, err := m.repo.queries(ctx).MailCreate(ctx, sqlc.MailCreateParams{
			Content:  mail.Content,
			SendTime: pgtype.Timestamptz{Valid: true, Time: mail.SendTime},
			Send:     mail.Send,
			Error:    pgtype.Text{Valid: mail.Error != "", String: mail.Error},
		})
		if err != nil {
			return fmt.Errorf("create mail %+v | %w", *mail, err)
		}

		mail.ID = int(id)

		mailEvents := make([]model.MailEvent, 0, len(eventIDs))
		for _, eventID := range eventIDs {
			mailEvents = append(mailEvents, model.MailEvent{
				MailID:  mail.ID,
				EventID: eventID,
			})
		}

		if err := m.mailEvent.CreateBatch(ctx, mailEvents); err != nil {
			return err
		}

		return nil
	})
}

func (m *Mail) Update(ctx context.Context, mail model.Mail, eventIDs []int) error {
	return m.repo.WithRollback(ctx, func(ctx context.Context) error {
		if err := m.repo.queries(ctx).MailUpdate(ctx, sqlc.MailUpdateParams{
			ID:       int32(mail.ID),
			Content:  mail.Content,
			SendTime: pgtype.Timestamptz{Valid: true, Time: mail.SendTime},
		}); err != nil {
			return fmt.Errorf("update mail %+v | %w", mail, err)
		}

		if err := m.mailEvent.DeleteByMail(ctx, mail.ID); err != nil {
			return err
		}

		mailEvents := make([]model.MailEvent, 0, len(eventIDs))
		for _, eventID := range eventIDs {
			mailEvents = append(mailEvents, model.MailEvent{
				MailID:  mail.ID,
				EventID: eventID,
			})
		}

		if err := m.mailEvent.CreateBatch(ctx, mailEvents); err != nil {
			return err
		}

		return nil
	})
}

func (m *Mail) Send(ctx context.Context, mailID int) error {
	if err := m.repo.queries(ctx).MailSend(ctx, int32(mailID)); err != nil {
		return fmt.Errorf("send mail %d | %w", mailID, err)
	}

	return nil
}

func (m *Mail) Error(ctx context.Context, mail model.Mail) error {
	if err := m.repo.queries(ctx).MailError(ctx, sqlc.MailErrorParams{
		ID:    int32(mail.ID),
		Error: pgtype.Text{Valid: true, String: mail.Error},
	}); err != nil {
		return fmt.Errorf("error mail %+v | %w", mail, err)
	}

	return nil
}
