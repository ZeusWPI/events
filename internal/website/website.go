// Package website scrapes the Zeus WPI website to get all event data
package website

import (
	"errors"

	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/pkg/config"
)

type Website struct {
	githubToken string

	eventRepo  repository.Event
	yearRepo   repository.Year
	boardRepo  repository.Board
	memberRepo repository.Member
}

func New(repo repository.Repository) (*Website, error) {
	github := config.GetDefaultString("website.github_token", "")
	if github == "" {
		return nil, errors.New("no github token set")
	}

	return &Website{
		githubToken: github,
		eventRepo:   *repo.NewEvent(),
		yearRepo:    *repo.NewYear(),
		boardRepo:   *repo.NewBoard(),
		memberRepo:  *repo.NewMember(),
	}, nil
}
