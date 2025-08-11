package mail

import (
	"context"

	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/task"
)

type Mail struct {
	repoMail repository.Mail
	task     *task.Manager
}

func New(repo repository.Repository, task *task.Manager) (*Mail, error) {
	mail := &Mail{
		repoMail: *repo.NewMail(),
		task:     task,
	}

	if err := mail.startup(context.Background()); err != nil {
		return nil, err
	}

	return mail, nil
}

func (m *Mail) startup(ctx context.Context) error {
	// Reschedule all mails
	mails, err := m.repoMail.GetUnsend(ctx)
	if err != nil {
		return err
	}

	for _, mail := range mails {
		if err := m.ScheduleMailAll(ctx, *mail, false); err != nil {
			return err
		}
	}

	return nil
}
