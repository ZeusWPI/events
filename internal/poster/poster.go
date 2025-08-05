package poster

import (
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/pkg/gitmate"
)

const SyncTask = "Poster sync"

type Client struct {
	gitmate gitmate.Client

	event  repository.Event
	poster repository.Poster
	year   repository.Year
}

func New(repo repository.Repository) (*Client, error) {
	gitmate, err := gitmate.New()
	if err != nil {
		return nil, err
	}

	return &Client{
		gitmate: *gitmate,
		event:   *repo.NewEvent(),
		poster:  *repo.NewPoster(),
		year:    *repo.NewYear(),
	}, nil
}
