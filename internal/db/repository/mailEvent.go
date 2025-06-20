package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/sqlc"
	"github.com/ZeusWPI/events/pkg/utils"
)

type MailEvent struct {
	repo Repository
}

func (r *Repository) NewMailEvent() *MailEvent {
	return &MailEvent{
		repo: *r,
	}
}

func (m *MailEvent) GetByEvents(ctx context.Context, events []model.Event) ([]*model.MailEvent, error) {
	mails, err := m.repo.queries(ctx).MailEventGetByEvents(ctx, utils.SliceMap(events, func(e model.Event) int32 { return int32(e.ID) }))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get all mail events by events %+v | %w", events, err)
	}

	return utils.SliceMap(mails, model.MailEventModel), nil
}

func (m *MailEvent) CreateBatch(ctx context.Context, mails []model.MailEvent) error {
	if err := m.repo.queries(ctx).MailEventCreateBatch(ctx, sqlc.MailEventCreateBatchParams{
		Column1: utils.SliceMap(mails, func(m model.MailEvent) int32 { return int32(m.MailID) }),
		Column2: utils.SliceMap(mails, func(m model.MailEvent) int32 { return int32(m.EventID) }),
	}); err != nil {
		return fmt.Errorf("create mail event %+v | %w", mails, err)
	}

	return nil
}

func (m *MailEvent) DeleteByMail(ctx context.Context, mailID int) error {
	if err := m.repo.queries(ctx).MailEventDeleteByMail(ctx, int32(mailID)); err != nil {
		return fmt.Errorf("delete mail event by mail %d | %w", mailID, err)
	}

	return nil
}
