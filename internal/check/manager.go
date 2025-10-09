package check

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/task"
	"github.com/ZeusWPI/events/pkg/config"
	"github.com/ZeusWPI/events/pkg/mattermost"
	"github.com/ZeusWPI/events/pkg/utils"
	"go.uber.org/zap"
)

const taskUID = "task-checkManager"

// Manager is the global single check manager instance
var Manager *manager

// manager manages all the checks
// After an application reboot each automatic check is set as inactive, aka they won't show up in the frontend
// So just like the task manager you need to reregister checks after an application reboot which will set them as active again
// Check statusses are saved in the databses however the caller is reponsible for updating statuses between reboots
// The UID should always be the same, otherwise the previous check results are lost (but still in the database)
// If you want to change the frontend display description then you should change the return value of the Description() function
type manager struct {
	checks []model.Check

	repoCheck repository.Check
	repoEvent repository.Event

	development bool
	mattermost  mattermost.Client
	channelID   string

	// The mutexes main purpose is to avoid any concurrent changes in the dabase to the automatic checks
	// Mainly for the syncDeadline function as well times function calls can overwrite statusses
	// The check manager is not used from the API so there won't be any noticable delay
	// because of the non concurrent db transactions.
	mu sync.Mutex
}

func newManager(repo repository.Repository) (*manager, error) {
	mClient, err := mattermost.New()
	if err != nil {
		return nil, err
	}

	manager := &manager{
		repoCheck:   *repo.NewCheck(),
		repoEvent:   *repo.NewEvent(),
		development: config.IsDev(),
		mattermost:  *mClient,
		channelID:   config.GetString("check.channel"),
	}

	if err := manager.repoCheck.SetInactiveAutomatic(context.Background()); err != nil {
		return nil, err
	}

	if err := task.Manager.AddRecurring(context.Background(), task.NewTask(
		taskUID,
		"Check status update",
		config.GetDefaultDuration("check.sync_s", 3*60*60),
		manager.syncDeadline,
	)); err != nil {
		return nil, err
	}

	return manager, nil
}

// Register a new check
// The id has to be unique
// If the id ever changes then the statusses of previous events will be lost
// If the status for an event is TODO and there is a deadline (!= NoDeadline)
// then it will automatically update the status to TODOLate when the deadline passed
// If NoDeadline is used as deadline then you're responsible for marking a check as late
func (m *manager) Register(ctx context.Context, newCheck Check) error {
	zap.S().Infof("Adding check: %s | deadline: %s", newCheck.Description(), newCheck.Deadline())

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, c := range m.checks {
		if c.UID == newCheck.UID() {
			return fmt.Errorf("check %s already exists", newCheck.UID())
		}
	}

	check, err := m.repoCheck.GetByUID(ctx, newCheck.UID())
	if err != nil {
		return err
	}
	if check != nil {
		// Pre-existing check
		// Update it
		check.Description = newCheck.Description()
		check.Deadline = newCheck.Deadline()
		check.Active = true
		if err := m.repoCheck.Update(ctx, *check); err != nil {
			return err
		}
	} else {
		// New check
		// Let's create it
		check = &model.Check{
			UID:         newCheck.UID(),
			Description: newCheck.Description(),
			Deadline:    newCheck.Deadline(),
			Active:      true,
			Type:        model.CheckAutomatic,
		}
		if err := m.repoCheck.Create(ctx, *check); err != nil {
			return err
		}
	}

	// Pre-populate the database with TODO and TODO late
	// If the check already exisited then this would only be for new events
	// If it's a new check then it's going to add an entry for every event
	checks, err := m.repoCheck.GetEventsByCheckUID(ctx, newCheck.UID())
	if err != nil {
		return err
	}
	events, err := m.repoEvent.GetAll(ctx)
	if err != nil {
		return err
	}

	toCreate := []model.Check{}
	for _, event := range events {
		if idx := slices.IndexFunc(checks, func(c *model.Check) bool { return c.EventID == event.ID }); idx != -1 {
			// There already exist a check for this event
			continue
		}

		status := model.CheckTODO
		if check.Deadline != NoDeadline && time.Now().Add(check.Deadline).After(event.StartTime) {
			status = model.CheckTODOLate
		}

		toCreate = append(toCreate, model.Check{
			UID:     check.UID,
			EventID: event.ID,
			Status:  status,
			Message: "",
		})
	}

	if err := m.repoCheck.CreateEventBatch(ctx, toCreate); err != nil {
		return err
	}

	m.checks = append(m.checks, *check)

	return nil
}

// Update let's you update the status of a check for an event
func (m *manager) Update(ctx context.Context, checkUID string, update Update) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := utils.SliceFind(m.checks, func(c model.Check) bool { return c.UID == checkUID }); !ok {
		return fmt.Errorf("check %s is not registered", checkUID)
	}

	check, err := m.repoCheck.GetByCheckUIDEvent(ctx, checkUID, update.EventID)
	if err != nil {
		return err
	}
	if check == nil {
		// This is only possible in the following rare cases
		//   - After a website event sync it got added to the db but the check NewEvent function failed
		// We could prevent it by periodically checking for new events in case an insert failed but this is way simpler
		check = &model.Check{
			UID:     checkUID,
			EventID: update.EventID,
			Status:  model.CheckTODO, // Doesn't matter, will immediatly be changed
		}

		if err := m.repoCheck.CreateEvent(ctx, check); err != nil {
			return err
		}
	}

	if check.Status == update.Status && check.Message == update.Message {
		// No change, no need to update it
		return nil
	}

	// Update the entry
	oldStatus := check.Status
	check.Status = update.Status
	check.Message = update.Message

	if err := m.repoCheck.UpdateEvent(ctx, *check); err != nil {
		return err
	}

	message := fmt.Sprintf("**Check Update**\n%s\n`%s` -> `%s`", check.Description, oldStatus, check.Status)
	if m.development {
		zap.S().Infof("Mock check update: \n%s", message)
	} else {
		if err := m.mattermost.SendMessage(ctx, mattermost.Message{
			ChannelID: m.channelID,
			Message:   message,
		}); err != nil {
			return err
		}
	}

	return nil
}

// NewEvent handles a new event being created
// It will create new check_events for each registered check
func (m *manager) NewEvent(ctx context.Context, event model.Event) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	checkEvents := make([]model.Check, 0, len(m.checks))
	for _, check := range m.checks {
		status := model.CheckTODO
		if check.Deadline != NoDeadline && time.Now().Add(check.Deadline).After(event.StartTime) {
			status = model.CheckTODOLate
		}

		checkEvents = append(checkEvents, model.Check{
			UID:     check.UID,
			EventID: event.ID,
			Status:  status,
		})
	}

	if err := m.repoCheck.CreateEventBatch(ctx, checkEvents); err != nil {
		return err
	}

	return nil
}

// syncDeadline will go over each check_event and update the status to TODOLate if the following conditions are met
//   - Status == TODO
//   - Type == "automatic"
//   - Deadline != NoDeadline
//   - Active == true
//   - Event startTime is in the future
func (m *manager) syncDeadline(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	events, err := m.repoEvent.GetFuture(ctx)
	if err != nil {
		return err
	}
	if events == nil {
		return nil
	}

	checks, err := m.repoCheck.GetByEvents(ctx, utils.SliceDereference(events))
	if err != nil {
		return err
	}
	for _, check := range checks {
		// Apply all the filters

		// Active == true is forced by the query
		// Event startTime in the future is forced by the events query

		// Status == TODO
		if check.Status != model.CheckTODO {
			continue
		}

		// Type == "automatic"
		if check.Type != model.CheckAutomatic {
			continue
		}

		// Deadline != NoDeadline
		if check.Deadline == NoDeadline {
			continue
		}

		// Deadline is in the future
		event, ok := utils.SliceFind(events, func(e *model.Event) bool { return e.ID == check.EventID })
		if !ok {
			return fmt.Errorf("no associated event with check %+v", *check)
		}
		if time.Now().Add(check.Deadline).Before(event.StartTime) {
			// Deadline is in the future, the board still has some time left
			// However if there's less than 1 day left send a message to mattermost to warn them
			if time.Now().Add(check.Deadline).Add(24 * time.Hour).After(event.StartTime) {
				if check.Mattermost {
					// We already warned them
					continue
				}

				message := fmt.Sprintf("**⚠️Check Deadline⚠️**\nCheck: `%s`\nEvent: `%s`\nTime left: `%s`", check.Description, event.Name, event.StartTime.Sub(time.Now().Add(check.Deadline)))
				if m.development {
					zap.S().Infof("Mock deadline warning: \n%s", message)
				} else {
					if err := m.mattermost.SendMessage(ctx, mattermost.Message{
						ChannelID: m.channelID,
						Message:   message,
					}); err != nil {
						return err
					}
				}

				if err := m.repoCheck.SendMattermost(ctx, check.ID); err != nil {
					return err
				}
			}
			continue
		}

		// All filters passes
		// Update the check
		check.Status = model.CheckTODOLate
		if err := m.repoCheck.UpdateEvent(ctx, *check); err != nil {
			return err
		}
	}

	return nil
}
