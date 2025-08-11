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
}

func (r *Repository) NewMail() *Mail {
	return &Mail{
		repo: *r,
	}
}

func (m *Mail) GetByID(ctx context.Context, mailID int) (*model.Mail, error) {
	mails, err := m.repo.queries(ctx).MailGetByID(ctx, int32(mailID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get mail by id %d | %w", mailID, err)
	}

	return model.MailEventsModel(mails)[0], nil
}

func (m *Mail) GetByYear(ctx context.Context, yearID int) ([]*model.Mail, error) {
	mails, err := m.repo.queries(ctx).MailGetByYear(ctx, int32(yearID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get mails by year %w", err)
	}

	return model.MailEventsModel(utils.SliceMap(mails, func(m sqlc.MailGetByYearRow) sqlc.MailGetByIDRow { return sqlc.MailGetByIDRow(m) })), nil
}

func (m *Mail) GetByEvents(ctx context.Context, events []model.Event) ([]*model.Mail, error) {
	mails, err := m.repo.queries(ctx).MailGetByEvents(ctx, utils.SliceMap(events, func(e model.Event) int32 { return int32(e.ID) }))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get mails by events %w", err)
	}

	return model.MailEventsModel(utils.SliceMap(mails, func(m sqlc.MailGetByEventsRow) sqlc.MailGetByIDRow { return sqlc.MailGetByIDRow(m) })), nil
}

func (m *Mail) GetUnsend(ctx context.Context) ([]*model.Mail, error) {
	mails, err := m.repo.queries(ctx).MailGetUnsend(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get unsend mails %w", err)
	}

	return model.MailEventsModel(utils.SliceMap(mails, func(m sqlc.MailGetUnsendRow) sqlc.MailGetByIDRow { return sqlc.MailGetByIDRow(m) })), nil
}

func (m *Mail) Create(ctx context.Context, mail *model.Mail) error {
	return m.repo.WithRollback(ctx, func(ctx context.Context) error {
		id, err := m.repo.queries(ctx).MailCreate(ctx, sqlc.MailCreateParams{
			YearID:   int32(mail.YearID),
			Title:    mail.Title,
			Content:  mail.Content,
			SendTime: pgtype.Timestamptz{Valid: true, Time: mail.SendTime},
			Send:     mail.Send,
			Error:    pgtype.Text{Valid: mail.Error != "", String: mail.Error},
		})
		if err != nil {
			return fmt.Errorf("create mail %+v | %w", *mail, err)
		}

		mail.ID = int(id)

		if len(mail.EventIDs) > 0 {
			if err := m.repo.queries(ctx).MailEventCreateBatch(ctx, sqlc.MailEventCreateBatchParams{
				Column1: utils.SliceRepeat(id, len(mail.EventIDs)),
				Column2: utils.SliceMap(mail.EventIDs, func(id int) int32 { return int32(id) }),
			}); err != nil {
				return fmt.Errorf("create mail events %+v | %w", *mail, err)
			}
		}

		return nil
	})
}

func (m *Mail) Update(ctx context.Context, mail model.Mail) error {
	return m.repo.WithRollback(ctx, func(ctx context.Context) error {
		if err := m.repo.queries(ctx).MailUpdate(ctx, sqlc.MailUpdateParams{
			ID:       int32(mail.ID),
			Title:    mail.Title,
			Content:  mail.Content,
			SendTime: pgtype.Timestamptz{Valid: true, Time: mail.SendTime},
		}); err != nil {
			return fmt.Errorf("update mail %+v | %w", mail, err)
		}

		if err := m.repo.queries(ctx).MailEventDeleteByMail(ctx, int32(mail.ID)); err != nil {
			return fmt.Errorf("update mail events (delete) %+v | %w", mail, err)
		}

		if len(mail.EventIDs) > 0 {
			if err := m.repo.queries(ctx).MailEventCreateBatch(ctx, sqlc.MailEventCreateBatchParams{
				Column1: utils.SliceRepeat(int32(mail.ID), len(mail.EventIDs)),
				Column2: utils.SliceMap(mail.EventIDs, func(id int) int32 { return int32(id) }),
			}); err != nil {
				return fmt.Errorf("update mail events (insert) %+v | %w", mail, err)
			}
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

func (m *Mail) Delete(ctx context.Context, mailID int) error {
	if err := m.repo.queries(ctx).MailDelete(ctx, int32(mailID)); err != nil {
		return fmt.Errorf("delete mail %d | %w", mailID, err)
	}

	return nil
}
