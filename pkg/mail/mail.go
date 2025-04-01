// Package mail provides functionality to send mails
package mail

import (
	"context"
	"strings"

	"github.com/ZeusWPI/events/pkg/config"
)

type mailer interface {
	send(context.Context, Content) error
}

// Content contains all data to send an email
type Content struct {
	Recipients []string
	Subject    string

	HTML string
	Data interface{}
}

// Client allows for sending mails using a specific mailer
type Client struct {
	mailer mailer
}

// NewClient constructs a new Client struct.
// It selects the mailers based on the selected environment.
// In dev environments it will show the email in the browser.
// In prod environments it will actually send the mail
func NewClient() *Client {
	env := config.GetDefaultString("app.env", "development")
	env = strings.ToLower(env)

	var mailer mailer

	if env == "development" {
		mailer = newBrowser()
	} else {
		mailer = newSMTP()
	}

	return &Client{
		mailer: mailer,
	}
}

// Send sends a new mail.
// As context it uses a default background context
func (c *Client) Send(content Content) error {
	return c.mailer.send(context.Background(), content)
}

// SendWithCtx sends a new mail and supports a context with a deadline
func (c *Client) SendWithCtx(ctx context.Context, content Content) error {
	return c.mailer.send(ctx, content)
}
