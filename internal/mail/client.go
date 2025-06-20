package mail

import (
	"context"
	"fmt"
	"html/template"
	"strings"

	"github.com/ZeusWPI/events/pkg/config"
	"github.com/wneessen/go-mail"
)

type mailer interface {
	mail(context.Context, content) error
}

type content struct {
	Recipients []string
	Subject    string

	HTML string
	Data interface{}
}

type client struct {
	mailer mailer
}

// newClient constructs a new Client struct.
// It selects the mailers based on the selected environment.
// In dev environments it will show the email in the browser.
// In prod environments it will actually send the mail
func newClient() *client {
	env := config.GetDefaultString("app.env", "development")
	env = strings.ToLower(env)

	var mailer mailer

	if env == "development" {
		mailer = newBrowser()
	} else {
		mailer = newSMTP()
	}

	return &client{
		mailer: mailer,
	}
}

func (c *client) mail(ctx context.Context, content content) error {
	return c.mailer.mail(ctx, content)
}

func newMail(content content) (*mail.Msg, error) {
	m := mail.NewMsg()

	if err := m.FromFormat("Events", "events@zeus.ugent.be"); err != nil {
		return nil, fmt.Errorf("unable to set mail from address %w", err)
	}
	if err := m.To(content.Recipients...); err != nil {
		return nil, fmt.Errorf("unable to set mail recipients %v | %w", content.Recipients, err)
	}

	m.Subject(content.Subject)
	m.SetDate()

	tmpl, err := template.New("htmltpl").Parse(content.HTML)
	if err != nil {
		return nil, fmt.Errorf("unable to parse html template %w", err)
	}

	if err := m.SetBodyHTMLTemplate(tmpl, content.Data); err != nil {
		return nil, fmt.Errorf("failed to add html template to mail body %w", err)
	}

	return m, nil
}
