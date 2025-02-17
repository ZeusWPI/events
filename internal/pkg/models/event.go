// Package models contains all internal data models
package models

import (
	"time"
)

// Event represents a Zeus WPI event
type Event struct {
	ID           int
	URL          string
	Name         string
	Description  string
	StartTime    time.Time
	EndTime      time.Time
	AcademicYear string
	Location     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
}
