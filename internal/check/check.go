package check

import (
	"context"

	"github.com/ZeusWPI/events/internal/db/model"
)

type Check interface {
	Description() string
	Status(ctx context.Context, events []model.Event) []CheckResult
}

type Status string

const (
	Finished   Status = "finished"
	Unfinished Status = "unfinished"
	Warning    Status = "warning"
)

// CheckResult is the result of a status query
type CheckResult struct {
	EventID int
	Status  Status
	Warning string
	Error   error
}

type Source string

const (
	Automatic Source = "automatic"
	Manual    Source = "manual"
)

// EventStatus contains all the info of a check
// Has some optional fields depending on the source
type EventStatus struct {
	ID          int // Only for tasks entered in the website
	EventID     int
	Description string
	Warning     string // Clarification why the status == warning (if applicable)
	Status      Status
	Error       error
	Source      Source
}
