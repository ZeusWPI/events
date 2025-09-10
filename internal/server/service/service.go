// Package service provides all business logic required by the api and converts between dto and models
package service

import (
	"context"

	"github.com/ZeusWPI/events/internal/announcement"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/mail"
	"github.com/ZeusWPI/events/internal/poster"
	"github.com/ZeusWPI/events/internal/website"
)

// Service is used to create specific services
type Service struct {
	repo repository.Repository

	mail          mail.Client
	announcements announcement.Client
	website       *website.Client
	poster        poster.Client
}

// New creates a new Service
func New(repo repository.Repository, mail mail.Client, website *website.Client, announcement announcement.Client, poster poster.Client) *Service {
	return &Service{
		repo:          repo,
		mail:          mail,
		website:       website,
		announcements: announcement,
		poster:        poster,
	}
}

func (s *Service) withRollback(ctx context.Context, fn func(context.Context) error) error {
	return s.repo.WithRollback(ctx, fn)
}
