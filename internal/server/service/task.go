package service

import (
	"context"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/server/dto"
	"github.com/ZeusWPI/events/pkg/utils"
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
	tasks, err := t.service.manager.Tasks()
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
	return t.service.manager.Run(id)
}
