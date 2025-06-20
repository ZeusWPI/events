package mail

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/task"
	"github.com/ZeusWPI/events/pkg/zauth"
)

const mailTask = "Mail send"

func (m *Mail) sendMailAll(ctx context.Context, mail model.Mail) error {
	if err := zauth.MailAll(ctx, mail.Content); err != nil {
		mail.Error = err.Error()
		if dbErr := m.repoMail.Error(ctx, mail); dbErr != nil {
			err = errors.Join(err, dbErr)
		}

		return err
	}

	if err := m.repoMail.Send(ctx, mail.ID); err != nil {
		return err
	}

	return nil
}

// If a mail is already scheduled then update needs to be set to true so that it cancels it first
func (m *Mail) ScheduleMailAll(ctx context.Context, mail model.Mail, update bool) error {
	name := fmt.Sprintf("%s %d", mailTask, mail.ID)

	if update {
		_ = m.task.RemoveOnceByName(name)
	}

	if mail.SendTime.Before(time.Now()) {
		mail.Error = "Mail send time was in the past"
		if err := m.repoMail.Error(ctx, mail); err != nil {
			return err
		}

		return nil
	}

	interval := time.Until(mail.SendTime)

	if err := m.task.AddOnce(task.NewTask(
		name,
		interval,
		func(ctx context.Context) error { return m.sendMailAll(ctx, mail) },
	)); err != nil {
		return err
	}

	return nil
}
