// Package gitmate interacts with the visueel gitmate repository
package gitmate

import (
	"errors"

	"github.com/ZeusWPI/events/pkg/config"
)

type Client struct {
	url   string
	token string
}

func New() (*Client, error) {
	token := config.GetDefaultString("gitmate.token", "")
	if token == "" {
		return nil, errors.New("no gitmate token set")
	}

	url := config.GetDefaultString("gitmate.url", "")
	if token == "" {
		return nil, errors.New("no gitmate url set")
	}

	return &Client{
		url:   url,
		token: token,
	}, nil
}
