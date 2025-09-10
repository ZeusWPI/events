package service

import (
	"context"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/server/dto"
	"github.com/ZeusWPI/events/internal/task"
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
	tasks, err := task.Manager.Tasks()
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}
	if tasks == nil {
		return []dto.Task{}, nil
	}

	return utils.SliceMap(tasks, dto.TaskDTO), nil
}

func (t *Task) GetHistory(ctx context.Context, filter dto.TaskFilter) ([]dto.TaskHistory, error) {
	tasks, err := t.task.GetFiltered(ctx, model.TaskFilter(filter))
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}
	if tasks == nil {
		return []dto.TaskHistory{}, nil
	}

	return utils.SliceMap(tasks, dto.TaskHistoryDTO), nil
}

func (t *Task) Start(taskUID string) error {
	return task.Manager.RunByUID(taskUID)
}

func (t *Task) Resolve(ctx context.Context, runID int) error {
	run, err := t.task.GetByRunID(ctx, runID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}
	if run == nil {
		return fiber.ErrNotFound
	}

	run.Result = model.Resolved

	if err := t.task.RunResolve(ctx, runID); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return nil
}
