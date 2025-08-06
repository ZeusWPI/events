package poster

import (
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/pkg/config"
	"github.com/ZeusWPI/events/pkg/gitmate"
)

const SyncTask = "Poster sync"

type Client struct {
	development bool
	gitmate     gitmate.Client

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
		development: config.GetDefaultString("app.env", "development") == "development",
		gitmate:     *gitmate,
		event:       *repo.NewEvent(),
		poster:      *repo.NewPoster(),
		year:        *repo.NewYear(),
	}, nil
}
