package poster

import (
	"context"

	"github.com/ZeusWPI/events/internal/check"
	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/pkg/utils"
)

type CheckBig struct {
	poster repository.Poster
}

func (c *Client) NewCheckBig() *CheckBig {
	return &CheckBig{
		poster: c.poster,
	}
}

// Interface compliance
var _ check.Check = (*CheckBig)(nil)

func (c *CheckBig) Description() string {
	return "Design a big event poster"
}

func (c *CheckBig) Status(ctx context.Context, events []model.Event) []check.CheckResult {
	statusses := make(map[int]check.CheckResult)
	for _, event := range events {
		statusses[event.ID] = check.CheckResult{
			EventID: event.ID,
			Status:  check.Unfinished,
			Error:   nil,
		}
	}

	posters, err := c.poster.GetByEvents(ctx, events)
	if err != nil {
		for k, v := range statusses {
			v.Error = err
			statusses[k] = v
		}

		return utils.MapValues(statusses)
	}

	for _, poster := range posters {
		if poster.SCC {
			continue
		}

		if status, ok := statusses[poster.EventID]; ok {
			status.Status = check.Finished
			statusses[poster.EventID] = status
		}
	}

	return utils.MapValues(statusses)
}

type CheckSCC struct {
	poster repository.Poster
}

func (c *Client) NewCheckSCC() *CheckSCC {
	return &CheckSCC{
		poster: c.poster,
	}
}

// Interface compliance
var _ check.Check = (*CheckSCC)(nil)

func (c *CheckSCC) Description() string {
	return "Design a scc event poster"
}

func (c *CheckSCC) Status(ctx context.Context, events []model.Event) []check.CheckResult {
	statusses := make(map[int]check.CheckResult)
	for _, event := range events {
		statusses[event.ID] = check.CheckResult{
			EventID: event.ID,
			Status:  check.Unfinished,
			Error:   nil,
		}
	}

	posters, err := c.poster.GetByEvents(ctx, events)
	if err != nil {
		for k, v := range statusses {
			v.Error = err
			statusses[k] = v
		}

		return utils.MapValues(statusses)
	}

	for _, poster := range posters {
		if !poster.SCC {
			continue
		}

		if status, ok := statusses[poster.EventID]; ok {
			status.Status = check.Finished
			statusses[poster.EventID] = status
		}
	}

	return utils.MapValues(statusses)
}
