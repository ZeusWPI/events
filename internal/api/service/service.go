// Package service provides all business logic required by the api and converts between dto and models
package service

import (
	"context"

	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/task"
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

func (s *Service) NewEvent() *Event {
	return newEvent(*s)
}

func (s *Service) NewOrganizer() *Organizer {
	return newOrganizer(*s)
}

func (s *Service) NewYear() *Year {
	return newYear(*s)
}

func (s *Service) NewTask() *Task {
	return newTask(*s)
}
