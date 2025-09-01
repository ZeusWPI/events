package service

import (
	"context"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/server/dto"
	"github.com/ZeusWPI/events/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Task struct {
	service Service

	task repository.Task
}

func (s *Service) NewTask() *Task {
	return &Task{
		service: *s,
		task:    *s.repo.NewTask(),
	}
}

func (t *Task) GetAll() ([]dto.Task, error) {
	tasks, err := t.service.task.Tasks()
	if err != nil {
		return nil, err
	}
	if tasks == nil {
		return []dto.Task{}, nil
	}

	return utils.SliceMap(tasks, dto.TaskDTO), nil
}

func (t *Task) GetHistory(ctx context.Context, filters dto.TaskHistoryFilter) ([]dto.TaskHistory, error) {
	tasks, err := t.task.GetFiltered(ctx, model.TaskFilter(filters))
	if err != nil {
		return nil, err
	}
	if tasks == nil {
		return []dto.TaskHistory{}, nil
	}

	return utils.SliceMap(tasks, dto.TaskHistoryDTO), nil
}

func (t *Task) Start(id int) error {
	return t.service.task.Run(id)
}

func (t *Task) Resolve(ctx context.Context, taskID int) error {
	task, err := t.task.Get(ctx, taskID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}
	if task == nil {
		return fiber.ErrNotFound
	}

	task.Result = model.Resolved

	if err := t.task.UpdateResult(ctx, *task); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return nil
}
