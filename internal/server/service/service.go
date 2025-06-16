// Package service provides all business logic required by the api and converts between dto and models
package service

import (
	"context"

	"github.com/ZeusWPI/events/internal/check"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/task"
	"github.com/ZeusWPI/events/internal/website"
)

// Service is used to create specific services
type Service struct {
	repo    repository.Repository
	check   check.Manager
	task    *task.Manager
	website website.Website
}

// New creates a new Service
func New(repo repository.Repository, check check.Manager, task *task.Manager, website website.Website) *Service {
	return &Service{
		repo:    repo,
		check:   check,
		task:    task,
		website: website,
	}
}

func (s *Service) withRollback(ctx context.Context, fn func(context.Context) error) error {
	return s.repo.WithRollback(ctx, fn)
}
