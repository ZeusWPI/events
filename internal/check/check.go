package check

import (
	"context"

	"github.com/ZeusWPI/events/internal/db/model"
)

type Check interface {
	Description() string
	Status(ctx context.Context, events []model.Event) []StatusResult
}

// The result of a status query
type StatusResult struct {
	EventID int
	Done    bool
	Error   error
}

type Source string

const (
	Automatic Source = "automatic"
	Manual    Source = "manual"
)

// Status contains all the info of a check
// Has some optional fields depending on the source
type Status struct {
	ID          int // Only for manual tasks
	EventID     int
	Description string
	Done        bool
	Error       error
	Source      Source
}
