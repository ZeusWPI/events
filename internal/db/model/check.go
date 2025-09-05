package model

import (
	"time"

	"github.com/ZeusWPI/events/internal/db/sqlc"
)

type CheckStatusEnum string

const (
	CheckSuccess CheckStatusEnum = "success"
	CheckFailed  CheckStatusEnum = "failed"
	CheckWarning CheckStatusEnum = "warning"
)

type Check struct {
	ID          int
	Description string
	Deadline    time.Duration
}

func CheckModel(c sqlc.Check) *Check {
	return &Check{
		ID:          int(c.ID),
		Description: c.Description,
		Deadline:    time.Duration(c.Deadline),
	}
}

type CheckStatus struct {
	ID      int
	EventID int
	CheckID int
	Status  CheckStatusEnum
	Message string
}

func CheckStatusModel(c sqlc.CheckStatus) *CheckStatus {
	message := ""
	if c.Message.Valid {
		message = c.Message.String
	}

	return &CheckStatus{
		ID:      int(c.ID),
		EventID: int(c.EventID),
		CheckID: int(c.CheckID),
		Status:  CheckStatusEnum(c.Status),
		Message: message,
	}
}

type CheckCustom struct {
	ID          int
	EventID     int
	Description string
	Status      CheckStatusEnum
	CreatorID   int
}

func CheckCustomModel(c sqlc.CheckCustom) *CheckCustom {
	return &CheckCustom{
		ID:          int(c.ID),
		EventID:     int(c.EventID),
		Description: c.Description,
		Status:      CheckStatusEnum(c.Status),
		CreatorID:   int(c.CreatorID),
	}
}
