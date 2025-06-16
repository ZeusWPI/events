package check

import (
	"context"
	"fmt"

	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/pkg/utils"
)

// Just like the task manager you need to reregister checks after an application reboot
type Manager struct {
	checks []Check

	repoCheck repository.Check
	repoEvent repository.Event
}

func NewManager(repo repository.Repository) *Manager {
	return &Manager{
		repoCheck: *repo.NewCheck(),
		repoEvent: *repo.NewEvent(),
	}
}

// Register a new check
// The name has to be unique
func (m *Manager) Register(check Check) error {
	name := check.Description()

	for _, c := range m.checks {
		if c.Description() == name {
			return fmt.Errorf("registered a duplicate check %s", name)
		}
	}

	m.checks = append(m.checks, check)

	return nil
}

func (m *Manager) Status(ctx context.Context, yearID int) (map[int][]Status, error) {
	eventsDB, err := m.repoEvent.GetByYearPopulated(ctx, yearID)
	if err != nil {
		return nil, err
	}
	if eventsDB == nil {
		return nil, nil
	}
	events := utils.SliceDereference(eventsDB)

	statusses := map[int][]Status{}

	// DB checks
	checks, err := m.repoCheck.GetByEvents(ctx, events)
	if err != nil {
		return nil, err
	}

	for _, check := range checks {
		status := Status{
			ID:          check.ID,
			EventID:     check.EventID,
			Description: check.Description,
			Done:        check.Done,
			Error:       nil,
			Source:      Manual,
		}

		statusses[check.EventID] = append(statusses[check.EventID], status)
	}

	// Registered checks
	for _, check := range m.checks {
		results := check.Status(ctx, events)
		name := check.Description()

		for _, result := range results {
			status := Status{
				ID:          0,
				EventID:     result.EventID,
				Description: name,
				Done:        result.Done,
				Error:       result.Error,
				Source:      Automatic,
			}

			statusses[result.EventID] = append(statusses[result.EventID], status)
		}
	}

	return statusses, nil
}
