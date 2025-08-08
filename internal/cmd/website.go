package cmd

import (
	"time"

	"github.com/ZeusWPI/events/internal/task"
	"github.com/ZeusWPI/events/internal/website"
	"github.com/ZeusWPI/events/pkg/config"
)

func Website(w *website.Client, t *task.Manager) error {
	// There's a webhook to trigger syncing but still run them periodically to be on the safe side
	if err := t.Add(task.NewTask(
		website.BoardTask,
		time.Duration(config.GetDefaultInt("website.boards_s", 86400))*time.Second,
		w.UpdateBoard,
	)); err != nil {
		return err
	}

	if err := t.Add(task.NewTask(
		website.EventTask,
		time.Duration(config.GetDefaultInt("website.events_s", 86400))*time.Second,
		w.UpdateEvent,
	)); err != nil {
		return err
	}

	return nil
}
