package mail

import (
	"context"
	"fmt"

	"github.com/ZeusWPI/events/pkg/config"
	"github.com/wneessen/go-mail"
)

type smtp struct {
	from     string
	fromAddr string

	host string
	port int
}

// Interface compliance
var _ mailer = (*smtp)(nil)

func newSMTP() *smtp {
	return &smtp{
		from:     config.GetDefaultString("mail.from", "Events"),
		fromAddr: config.GetDefaultString("mail.from_address", "events@zeus.ugent.be"),
		host:     config.GetDefaultString("mail.host", ""),
		port:     config.GetDefaultInt("mail.port", 25),
	}
}

func (s *smtp) mail(ctx context.Context, content content) error {
	m, err := newMail(content)
	if err != nil {
		return err
	}

	client, err := mail.NewClient(s.host, mail.WithPort(s.port))
	if err != nil {
		return fmt.Errorf("failed to create a new mail client %w", err)
	}

	if err := client.DialAndSendWithContext(ctx, m); err != nil {
		return fmt.Errorf("failed to send mail %w", err)
	}

	return nil
}
