package cmd

import (
	"time"

	"github.com/ZeusWPI/events/internal/poster"
	"github.com/ZeusWPI/events/internal/task"
	"github.com/ZeusWPI/events/pkg/config"
)

func Poster(p poster.Client, t *task.Manager) error {
	if err := t.Add(task.NewTask(
		poster.SyncTask,
		time.Duration(config.GetDefaultInt("poster.sync_s", 86400))*time.Second,
		p.Sync,
	)); err != nil {
		return err
	}

	return nil
}
