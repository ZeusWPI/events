package cmd

import (
	"time"

	"github.com/ZeusWPI/events/internal/check"
	"github.com/ZeusWPI/events/internal/dsa"
	"github.com/ZeusWPI/events/internal/task"
	"github.com/ZeusWPI/events/pkg/config"
)

func DSA(d *dsa.DSA, t *task.Manager, c *check.Manager) error {
	if err := t.Add(task.NewTask(
		dsa.ActivitiesTask,
		time.Duration(config.GetDefaultInt("dsa.get_activities_s", 3600))*time.Second,
		d.UpdateActivities,
	)); err != nil {
		return err
	}

	if err := c.Register(d); err != nil {
		return err
	}

	if err := t.Add(task.NewTask(
		dsa.CreateActivitesTask,
		time.Duration(config.GetDefaultInt("dsa.create_activities_s", 7200))*time.Second,
		d.CreateActivities,
	)); err != nil {
		return err
	}

	return nil
}
