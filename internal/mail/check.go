package mail

import (
	"context"

	"github.com/ZeusWPI/events/internal/check"
	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/pkg/utils"
)

type CheckMail struct {
	repoMail repository.Mail
}

func (m *Mail) NewCheckMail() *CheckMail {
	return &CheckMail{
		repoMail: m.repoMail,
	}
}

var _ check.Check = (*CheckMail)(nil)

func (c *CheckMail) Description() string {
	return "Cover the event in an mail"
}

func (c *CheckMail) Status(ctx context.Context, events []model.Event) []check.CheckResult {
	statusses := make(map[int]check.CheckResult)
	for _, event := range events {
		statusses[event.ID] = check.CheckResult{
			EventID: event.ID,
			Status:  check.Unfinished,
			Error:   nil,
		}
	}

	mails, err := c.repoMail.GetByEvents(ctx, events)
	if err != nil {
		for k, v := range statusses {
			v.Error = err
			statusses[k] = v
		}

		return utils.MapValues(statusses)
	}

	for _, mail := range mails {
		for _, eventID := range mail.EventIDs {
			if status, ok := statusses[eventID]; ok {
				status.Status = check.Finished
				statusses[eventID] = status
			}
		}
	}

	return utils.MapValues(statusses)
}
