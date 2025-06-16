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
	if err := w.service.manager.AddOnce(task.NewTask(website.BoardTask, task.Now, w.service.website.UpdateBoard)); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	if err := w.service.manager.AddOnce(task.NewTask(website.EventTask, task.Now, w.service.website.UpdateEvent)); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return nil
}
