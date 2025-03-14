// Package cmd contains all internal commands to start various tasks
package cmd

import (
	"time"

	"github.com/ZeusWPI/events/internal/pkg/website"
	"github.com/ZeusWPI/events/pkg/config"
)

// RunWebsitePeriodic starts all the periodic background tasks for the website package
func RunWebsitePeriodic(w *website.Website) {
	// Making sure fetching the years is run at least once before the events are fetched.
	_ = w.UpdateAllYears()

	yearsTask := &periodicTask{
		name:     "update years",
		interval: time.Duration(config.GetDefaultInt("website.years_s", 86400)) * time.Second,
		done:     make(chan bool),
		task:     func() error { return w.UpdateAllYears() },
	}
	go yearsTask.run()

	eventsTask := &periodicTask{
		name:     "update events",
		interval: time.Duration(config.GetDefaultInt("website.events_s", 3600)) * time.Second,
		done:     make(chan bool),
		task:     func() error { return w.UpdateAllEvents() },
	}
	go eventsTask.run()

	boardsTask := &periodicTask{
		name:     "update board members",
		interval: time.Duration(config.GetDefaultInt("website.boards_s", 86400)) * time.Second,
		done:     make(chan bool),
		task:     func() error { return w.UpdateAllBoards() },
	}
	go boardsTask.run()
}
