// Package cmd contains all internal commands to start various tasks
package cmd

import (
	"context"
	"time"

	"github.com/ZeusWPI/events/internal/task"
	"github.com/ZeusWPI/events/internal/website"
	"github.com/ZeusWPI/events/pkg/config"
)

// Website starts all background tasks for the website pkg
func Website(m *task.Manager, w website.Website) error {
	// Add fetching years
	if err := m.Add(task.NewTask(
		website.YearTask,
		time.Duration(config.GetDefaultInt("website.years_s", 86400))*time.Second,
		func(_ context.Context) error { return w.UpdateAllYears() },
	)); err != nil {
		return err
	}

	// Add fetching events
	if err := m.Add(task.NewTask(
		website.EventTask,
		time.Duration(config.GetDefaultInt("website.events_s", 3600))*time.Second,
		func(_ context.Context) error { return w.UpdateAllEvents() },
	)); err != nil {
		return err
	}

	// Add fetching board members
	if err := m.Add(task.NewTask(
		website.BoardTask,
		time.Duration(config.GetDefaultInt("website.boards_s", 86400))*time.Second,
		func(_ context.Context) error { return w.UpdateAllBoards() },
	)); err != nil {
		return err
	}

	return nil
}
