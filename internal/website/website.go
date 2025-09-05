// Package website scrapes the Zeus WPI website to get all event data
package website

import (
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/dsa"
	"github.com/ZeusWPI/events/pkg/github"
)

type Client struct {
	github     *github.Client
	eventRepo  repository.Event
	yearRepo   repository.Year
	boardRepo  repository.Board
	memberRepo repository.Member
	dsa        dsa.DSA
}

func New(repo repository.Repository) (*Client, error) {
	github, err := github.New()
	if err != nil {
		return nil, err
	}

	dsa, err := dsa.New(repo)
	if err != nil {
		return nil, err
	}

	return &Client{
		github:     github,
		eventRepo:  *repo.NewEvent(),
		yearRepo:   *repo.NewYear(),
		boardRepo:  *repo.NewBoard(),
		memberRepo: *repo.NewMember(),
		dsa:        *dsa,
	}, nil
}
