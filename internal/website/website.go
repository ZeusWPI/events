// Package website scrapes the Zeus WPI website to get all event data
package website

import (
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/pkg/github"
)

type Website struct {
	github     github.Client
	eventRepo  repository.Event
	yearRepo   repository.Year
	boardRepo  repository.Board
	memberRepo repository.Member
}

func New(repo repository.Repository) (*Website, error) {
	github, err := github.New()
	if err != nil {
		return nil, err
	}

	return &Website{
		github:     *github,
		eventRepo:  *repo.NewEvent(),
		yearRepo:   *repo.NewYear(),
		boardRepo:  *repo.NewBoard(),
		memberRepo: *repo.NewMember(),
	}, nil
}
