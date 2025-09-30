package mail

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/task"
	"github.com/ZeusWPI/events/pkg/zauth"
	"go.uber.org/zap"
)

func getTaskUID(mail model.Mail) string {
	return fmt.Sprintf("%s-%d", taskUID, mail.ID)
}

func getTaskName(mail model.Mail) string {
	return fmt.Sprintf("Send mail %d", mail.ID)
}

func (c *Client) sendMailAll(ctx context.Context, mail model.Mail) error {
	if c.development {
		// Mock the request in development
		zap.S().Infof("Mock mail: %+v", mail)
	} else {
		if err := zauth.MailAll(ctx, mail.Title, mail.Content); err != nil {
			mail.Error = err.Error()
			if dbErr := c.repoMail.Update(ctx, mail); dbErr != nil {
				err = errors.Join(err, dbErr)
			}

			return err
		}
	}

	if err := c.repoMail.Send(ctx, mail.ID); err != nil {
		return err
	}

	return nil
}

func (c *Client) ScheduleMailAll(ctx context.Context, mail model.Mail) error {
	name := getTaskName(mail)

	// Added as a failsafe but should be checked by the caller
	if mail.SendTime.Before(time.Now()) {
		mail.Error = "Mail send time is in the past"
		if err := c.repoMail.Update(ctx, mail); err != nil {
			return err
		}

		return nil
	}

	interval := time.Until(mail.SendTime)

	if err := task.Manager.AddOnce(task.NewTask(
		getTaskUID(mail),
		name,
		interval,
		func(ctx context.Context) error { return c.sendMailAll(ctx, mail) },
	)); err != nil {
		return err
	}

	return nil
}
