// Package website scrapes the Zeus WPI website to get all event data
package website

import (
	"context"

	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/dsa"
	"github.com/ZeusWPI/events/internal/task"
	"github.com/ZeusWPI/events/pkg/config"
	"github.com/ZeusWPI/events/pkg/github"
)

const (
	TaskEventsUID = "task-website-events"
	TaskBoardUID  = "task-website-board"
)

type Client struct {
	dsa    dsa.Client
	github *github.Client

	eventRepo  repository.Event
	yearRepo   repository.Year
	boardRepo  repository.Board
	memberRepo repository.Member
}

func New(repo repository.Repository, dsa dsa.Client) (*Client, error) {
	github, err := github.New()
	if err != nil {
		return nil, err
	}

	client := &Client{
		dsa:        dsa,
		github:     github,
		eventRepo:  *repo.NewEvent(),
		yearRepo:   *repo.NewYear(),
		boardRepo:  *repo.NewBoard(),
		memberRepo: *repo.NewMember(),
	}

	// Register tasks
	if err := task.Manager.AddRecurring(context.Background(), task.NewTask(
		TaskEventsUID,
		"Syncronize events",
		config.GetDefaultDuration("website.events_s", 24*60*60),
		client.SyncEvents,
	)); err != nil {
		return nil, err
	}

	if err := task.Manager.AddRecurring(context.Background(), task.NewTask(
		TaskBoardUID,
		"Syncronize boards",
		config.GetDefaultDuration("website.board_s", 24*60*60),
		client.SyncBoard,
	)); err != nil {
		return nil, err
	}

	return client, err
}
