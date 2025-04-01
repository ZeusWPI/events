package mail

import (
	"fmt"
	"html/template"

	"github.com/wneessen/go-mail"
)

func newMail(content Content) (*mail.Msg, error) {
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
