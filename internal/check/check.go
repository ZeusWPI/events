// Package check provides an interface to register event checks
package check

import (
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
)

// Init initiates the global check manager instance
func Init(repo repository.Repository) error {
	manager, err := newManager(repo)
	if err != nil {
		return err
	}

	Manager = manager

	return nil
}

// NoDeadline can be used if the check shouldn't have a deadline
// More specifically, it will prevent the manager from automatically changing
// the status to TODOLate when the deadline passes
const NoDeadline = time.Duration(0)

// Check is the interface to which a check should ahere to
// It exposes all the information needed to manage results
// You can manually implement all methods are make use of the `NewCheck` function
type Check interface {
	// UID is an unique string to identify the check
	// Results are linked with the uid
	// If the id changes then the previous results are gone
	UID() string
	// Description is an user friendly explanation of the check
	// You can change this as much as you like
	Description() string
	// Duration before the start of an event that it needs to be finished.
	// A duration of 0 means that it doesn't matter (NoDeadline).
	Deadline() time.Duration
}

// Update contains the information to update a check status
type Update struct {
	Status  model.CheckStatus
	Message string // Optional, let's you provide more info about the status
	EventID int
}

type internalCheck struct {
	uid         string
	description string
	deadline    time.Duration
}

var _ Check = (*internalCheck)(nil)

func NewCheck(uid string, description string, deadline time.Duration) Check {
	return &internalCheck{
		uid:         uid,
		description: description,
		deadline:    deadline,
	}
}

func (t *internalCheck) UID() string {
	return t.uid
}

func (t *internalCheck) Description() string {
	return t.description
}

func (t *internalCheck) Deadline() time.Duration {
	return t.deadline
}
