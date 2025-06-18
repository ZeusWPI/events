package cmd

import (
	"github.com/ZeusWPI/events/internal/check"
	"github.com/ZeusWPI/events/internal/mattermost"
)

func Mattermost(m *mattermost.Mattermost, c *check.Manager) error {
	if err := c.Register(m.Announcement); err != nil {
		return err
	}

	return nil
}
