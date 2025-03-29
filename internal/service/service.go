// Package service provides all business logic required by the api and converts between dto and models
package service

import (
	"context"

	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/pkg/task"
)

// Service is used to create specific services
type Service struct {
	repo    repository.Repository
	manager *task.Manager
}

// New creates a new Service
func New(repo repository.Repository, manager *task.Manager) *Service {
	return &Service{
		repo:    repo,
		manager: manager,
	}
}

func (s *Service) withRollback(ctx context.Context, fn func(context.Context) error) error {
	return s.repo.WithRollback(ctx, fn)
}

// NewEvent creates a new Event service
func (s *Service) NewEvent() Event {
	return &eventService{service: *s, manager: s.manager, board: s.repo.NewBoard(), event: s.repo.NewEvent(), organizer: s.repo.NewOrganizer()}
}

// NewOrganizer creates a new Organizer service
func (s *Service) NewOrganizer() Organizer {
	return &organizerService{service: *s, board: s.repo.NewBoard(), member: s.repo.NewMember(), year: s.repo.NewYear()}
}

// NewYear creates a new Year service
func (s *Service) NewYear() Year {
	return &yearService{service: *s, year: s.repo.NewYear()}
}

// NewTask creates a new Task service
func (s *Service) NewTask() Task {
	return &taskService{service: *s, manager: s.manager, task: s.repo.NewTask()}
}
