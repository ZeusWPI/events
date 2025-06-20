package cmd

import (
	"github.com/ZeusWPI/events/internal/check"
	"github.com/ZeusWPI/events/internal/mail"
)

func Mail(m *mail.Mail, c *check.Manager) error {
	if err := c.Register(m.NewCheckMail()); err != nil {
		return err
	}

	return nil
}
