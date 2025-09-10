// Package gitmate interacts with the visueel gitmate repository
package gitmate

import (
	"errors"

	"github.com/ZeusWPI/events/pkg/config"
)

// Change this url if you for some reason have a local instance
const url = "https://git.zeus.gent/api/v1/repos/ZeusWPI/visueel"

type Client struct {
	token string
}

func New() (*Client, error) {
	token := config.GetDefaultString("gitmate.token", "")
	if token == "" {
		return nil, errors.New("no gitmate token set")
	}

	return &Client{
		token: token,
	}, nil
}
