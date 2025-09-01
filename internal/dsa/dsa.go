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
	development  bool
	dsaURL       string
	dsaKey       string
	abbreviation string

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

	abbreviation := config.GetDefaultString("dsa.assoc_abbrev", "")
	if abbreviation == "" {
		return nil, errors.New("no association abbreviation set")
	}

	return &DSA{
		development:  config.GetDefaultString("app.env", "development") == "development",
		dsaURL:       url,
		dsaKey:       dsaKey,
		abbreviation: abbreviation,
		repoDSA:      *repo.NewDSA(),
		repoEvent:    *repo.NewEvent(),
		repoYear:     *repo.NewYear(),
	}, nil
}

// Interface compliance
var _ check.Check = (*DSA)(nil)

func (d *DSA) Description() string {
	return "Add event to the DSA website"
}

func (d *DSA) Status(ctx context.Context, events []model.Event) []check.CheckResult {
	statusses := make(map[int]check.CheckResult)
	for _, event := range events {
		statusses[event.ID] = check.CheckResult{
			EventID: event.ID,
			Status:  check.Unfinished,
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
			if dsa.Deleted {
				status.Status = check.Warning
				status.Warning = "DSA activity was manually deleted on the dsa website."
			} else if dsa.DsaID != 0 {
				status.Status = check.Finished
			}
			statusses[dsa.EventID] = status
		}
	}

	return utils.MapValues(statusses)
}
