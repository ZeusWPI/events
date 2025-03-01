// Package website scrapes the Zeus WPI website to get all event data
package website

import (
	"github.com/ZeusWPI/events/internal/pkg/db/repository"
)

// Warning: This package contains a lot of webscraping
// Webscraping results in ugly code

// Website represents the ZeusWPI website and all related functions
type Website struct {
	eventRepo  repository.Event
	yearRepo   repository.AcademicYear
	boardRepo  repository.Board
	memberRepo repository.Member
}

// New creates a new website instance
func New(repo repository.Repository) *Website {
	return &Website{
		eventRepo:  repo.NewEvent(),
		yearRepo:   repo.NewAcademicYear(),
		boardRepo:  repo.NewBoard(),
		memberRepo: repo.NewMember(),
	}
}
