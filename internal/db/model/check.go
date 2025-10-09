package model

import (
	"time"

	"github.com/ZeusWPI/events/internal/db/sqlc"
)

type CheckStatus string

const (
	CheckDone     CheckStatus = "done"      // Check is furfilled
	CheckDoneLate CheckStatus = "done_late" // Check is furfilled too late but without too many consequences
	CheckTODO     CheckStatus = "todo"      // Still needs to be done
	CheckTODOLate CheckStatus = "todo_late" // Still needs to be done but already too late
	CheckWarning  CheckStatus = "warning"   // Something is up with the check
)

type CheckType string

const (
	CheckManual    CheckType = "manual"
	CheckAutomatic CheckType = "automatic"
)

type Check struct {
	// Check event fields
	ID         int // ID of the check event
	EventID    int
	Status     CheckStatus
	Message    string
	Mattermost bool
	UpdatedAt  time.Time

	// Check fields
	UID         string // Identifier of the check
	Description string
	Deadline    time.Duration
	Active      bool
	Type        CheckType
	CreatorID   int
}

func CheckModel(check sqlc.Check, checkEvent sqlc.CheckEvent) *Check {
	message := ""
	if checkEvent.Message.Valid {
		message = checkEvent.Message.String
	}

	deadline := time.Duration(0)
	if check.Deadline.Valid {
		deadline = time.Duration(check.Deadline.Int64)
	}

	creatorID := 0
	if check.CreatorID.Valid {
		creatorID = int(check.CreatorID.Int32)
	}

	return &Check{
		ID:         int(checkEvent.ID),
		EventID:    int(checkEvent.EventID),
		Status:     CheckStatus(checkEvent.Status),
		Message:    message,
		Mattermost: checkEvent.Mattermost,
		UpdatedAt:  checkEvent.UpdatedAt.Time,

		UID:         check.Uid,
		Description: check.Description,
		Deadline:    deadline,
		Active:      check.Active,
		Type:        CheckType(check.Type),
		CreatorID:   creatorID,
	}
}
