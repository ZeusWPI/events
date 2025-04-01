package mail

import (
	"context"
	"fmt"
	"os"

	pkgBrowser "github.com/pkg/browser"
)

type browser struct{}

// Interface compliance
var _ mailer = (*browser)(nil)

func newBrowser() *browser {
	return &browser{}
}

func (b *browser) send(_ context.Context, content Content) error {
	m, err := newMail(content)
	if err != nil {
		return err
	}

	f, err := os.CreateTemp("", "mail.html")
	defer func() {
		_ = f.Close()
	}()

	if err != nil {
		return fmt.Errorf("unable to create temporarily mail file %w", err)
	}

	if _, err := m.WriteTo(f); err != nil {
		return fmt.Errorf("unable to write to temp file %w", err)
	}

	if err := pkgBrowser.OpenFile(f.Name()); err != nil {
		return fmt.Errorf("unable to open file in browser %s | %w", f.Name(), err)
	}

	return nil
}
