package dsa

import (
	"context"
	"errors"

	"github.com/ZeusWPI/events/internal/check"
	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/pkg/config"
	"github.com/ZeusWPI/events/pkg/utils"
)

type DSA struct {
	dsaURL string
	dsaKey string

	repoDSA   repository.DSA
	repoEvent repository.Event
	repoYear  repository.Year
}

func New(repo repository.Repository) (*DSA, error) {
	url := config.GetDefaultString("dsa.url", "")
	if url == "" {
		return nil, errors.New("no dsa url link set")
	}
	dsaKey := config.GetDefaultString("dsa.key", "")
	if dsaKey == "" {
		return nil, errors.New("no dsa api key set")
	}

	return &DSA{
		dsaURL:    url,
		dsaKey:    dsaKey,
		repoDSA:   *repo.NewDSA(),
		repoEvent: *repo.NewEvent(),
		repoYear:  *repo.NewYear(),
	}, nil
}

// Interface compliance
var _ check.Check = (*DSA)(nil)

func (d *DSA) Description() string {
	return "Event added to the DSA website"
}

func (d *DSA) Status(ctx context.Context, events []model.Event) []check.StatusResult {
	statusses := make(map[int]check.StatusResult)
	for _, event := range events {
		statusses[event.ID] = check.StatusResult{
			EventID: event.ID,
			Done:    false,
			Error:   nil,
		}
	}

	dsas, err := d.repoDSA.GetByEvents(ctx, events)
	if err != nil {
		for k, v := range statusses {
			v.Error = err
			statusses[k] = v
		}

		return utils.MapValues(statusses)
	}

	for _, dsa := range dsas {
		if status, ok := statusses[dsa.EventID]; ok {
			status.Done = dsa.Entry
			statusses[dsa.EventID] = status
		}
	}

	return utils.MapValues(statusses)
}
