package cmd

import (
	"time"

	"go.uber.org/zap"
)

type periodicTask struct {
	name     string
	interval time.Duration
	done     chan bool
	task     func() error
}

func (t *periodicTask) run() {
	zap.S().Infof("Starting periodic task %s with interval %s", t.name, t.interval)

	// Execute on startup
	if err := t.task(); err != nil {
		zap.S().Errorf("Periodic: Error during task execution %s | %v", t.name, err)
	}

	// Start periodic
	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()

	for {
		select {
		case <-t.done:
			zap.S().Warnf("Periodic: Stopping task %s", t.name)
			return

		case <-ticker.C:
			if err := t.task(); err != nil {
				zap.S().Errorf("Periodic: Error during task execution %s | %v", t.name, err)
			}
		}
	}
}
