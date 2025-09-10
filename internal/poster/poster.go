package poster

import (
	"context"
	"fmt"
	"time"

	"github.com/ZeusWPI/events/internal/check"
	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/task"
	"github.com/ZeusWPI/events/pkg/config"
	"github.com/ZeusWPI/events/pkg/gitmate"
)

const (
	checkBigUID = "check-poster-big"
	checkSCCUID = "check-poster-scc"
	TaskUID     = "task-poster"
)

type Client struct {
	development bool
	gitmate     gitmate.Client

	event  repository.Event
	poster repository.Poster
	year   repository.Year

	deadlineBig time.Duration
	deadlineSCC time.Duration
}

func New(repo repository.Repository) (*Client, error) {
	gitmate, err := gitmate.New()
	if err != nil {
		return nil, err
	}

	client := &Client{
		development: config.IsDev(),
		gitmate:     *gitmate,
		event:       *repo.NewEvent(),
		poster:      *repo.NewPoster(),
		year:        *repo.NewYear(),
		deadlineBig: config.GetDefaultDuration("poster.big_deadline_s", 0),
		deadlineSCC: config.GetDefaultDuration("poster.scc_deadline_s", 0),
	}

	// Register task
	if err := task.Manager.AddRecurring(context.Background(), task.NewTask(
		TaskUID,
		"Syncronize posters",
		config.GetDefaultDuration("poster.sync_s", 60*60),
		client.Sync,
	)); err != nil {
		return nil, err
	}

	// Register checks
	if err := check.Manager.Register(context.Background(), check.NewCheck(
		checkBigUID,
		"Add a big poster",
		client.deadlineBig,
	)); err != nil {
		return nil, err
	}
	if err := check.Manager.Register(context.Background(), check.NewCheck(
		checkSCCUID,
		"Add a scc poster",
		client.deadlineSCC,
	)); err != nil {
		return nil, err
	}

	if err := client.syncCheck(context.Background()); err != nil {
		return nil, err
	}

	return client, nil
}

// SyncCheck synchronizes the checks with the current db situation
// It's necessary to run at least once because of the new check system
// Should be removed in a later release
// TODO: remove
func (c *Client) syncCheck(ctx context.Context) error {
	posters, err := c.poster.GetAll(ctx)
	if err != nil {
		return err
	}

	for _, poster := range posters {
		event, err := c.event.GetByID(ctx, poster.EventID)
		if err != nil {
			return err
		}

		uid := checkBigUID
		if poster.SCC {
			uid = checkSCCUID
		}

		if err := check.Manager.Update(ctx, uid, check.Update{
			Status:  model.CheckDone,
			Message: "",
			EventID: event.ID,
		}); err != nil {
			return err
		}
	}

	return nil
}

// Create handles a poster being created
func (c *Client) Create(ctx context.Context, poster model.Poster) error {
	event, err := c.event.GetByID(ctx, poster.EventID)
	if err != nil {
		return err
	}
	if event == nil {
		return fmt.Errorf("no event found for poster %+v", poster)
	}

	deadline := c.deadlineBig
	uid := checkBigUID

	if poster.SCC {
		deadline = c.deadlineSCC
		uid = checkSCCUID
	}

	status := model.CheckDone
	if time.Now().Add(deadline).After(event.StartTime) {
		status = model.CheckDoneLate
	}

	if err := check.Manager.Update(ctx, uid, check.Update{
		Status:  status,
		EventID: event.ID,
	}); err != nil {
		return err
	}

	return nil
}

// Delete handles a poster being deleted
func (c *Client) Delete(ctx context.Context, poster model.Poster) error {
	event, err := c.event.GetByID(ctx, poster.EventID)
	if err != nil {
		return err
	}
	if event == nil {
		return fmt.Errorf("no event found for poster %+v", poster)
	}

	deadline := c.deadlineBig
	uid := checkBigUID
	if poster.SCC {
		deadline = c.deadlineSCC
		uid = checkSCCUID
	}

	status := model.CheckTODO
	if time.Now().Add(deadline).After(event.StartTime) {
		status = model.CheckTODOLate
	}

	if err := check.Manager.Update(ctx, uid, check.Update{
		Status:  status,
		EventID: event.ID,
	}); err != nil {
		return err
	}

	return nil
}
