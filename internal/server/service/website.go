package service

import (
	"github.com/ZeusWPI/events/internal/task"
	"github.com/ZeusWPI/events/internal/website"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Website struct {
	service Service
}

func (s *Service) NewWebsite() *Website {
	return &Website{
		service: *s,
	}
}

func (w *Website) Sync() error {
	// The task manager runs everything in the background
	// The returned error is the status for adding it to the task manager
	// The result of the task itself if logged by the task manager
	if err := task.Manager.AddOnce(task.NewTask(website.TaskBoardUID, "Updating boards due to repository update", task.Now, w.service.website.SyncBoard)); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	if err := task.Manager.AddOnce(task.NewTask(website.TaskEventsUID, "Updating events due to repository update", task.Now, w.service.website.SyncEvents)); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return nil
}
