package check

import (
	"context"
	"fmt"

	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/pkg/utils"
	"go.uber.org/zap"
)

// Manager manages all the checks
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
	zap.S().Infof("Adding check %s", check.Description())

	name := check.Description()

	for _, c := range m.checks {
		if c.Description() == name {
			return fmt.Errorf("registered a duplicate check %s", name)
		}
	}

	m.checks = append(m.checks, check)

	return nil
}

func (m *Manager) Status(ctx context.Context, yearID int) (map[int][]EventStatus, error) {
	eventsDB, err := m.repoEvent.GetByYearPopulated(ctx, yearID)
	if err != nil {
		return nil, err
	}
	if eventsDB == nil {
		return nil, nil
	}
	events := utils.SliceDereference(eventsDB)

	statusses := map[int][]EventStatus{}

	// DB checks
	checks, err := m.repoCheck.GetByEvents(ctx, events)
	if err != nil {
		return nil, err
	}

	for _, check := range checks {
		status := Unfinished
		if check.Done {
			status = Finished
		}

		eventStatus := EventStatus{
			ID:          check.ID,
			EventID:     check.EventID,
			Description: check.Description,
			Status:      status,
			Error:       nil,
			Source:      Manual,
		}

		statusses[check.EventID] = append(statusses[check.EventID], eventStatus)
	}
	// Registered checks
	for _, check := range m.checks {
		results := check.Status(ctx, events)
		name := check.Description()

		for _, result := range results {
			status := EventStatus{
				ID:          0,
				EventID:     result.EventID,
				Description: name,
				Warning:     result.Warning,
				Status:      result.Status,
				Error:       result.Error,
				Source:      Automatic,
			}

			statusses[result.EventID] = append(statusses[result.EventID], status)
		}
	}

	return statusses, nil
}
