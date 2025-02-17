// Package website scrapes the Zeus WPI website to get all event data
package website

import (
	"github.com/ZeusWPI/events/internal/pkg/db/repository"
)

// Website represents the ZeusWPI website and all related functions
type Website struct {
	eventRepo repository.Event
}

// New creates a new website instance
func New(repo repository.Repository) *Website {
	return &Website{
		eventRepo: repo.NewEvent(),
	}
}
