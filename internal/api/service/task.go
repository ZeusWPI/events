package service

import (
	"context"

	"github.com/ZeusWPI/events/internal/api/dto"
	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/task"
	"github.com/ZeusWPI/events/pkg/util"
)

// Task represents all business logic regarding tasks
type Task interface {
	GetAll() ([]dto.Task, error)
	GetHistory(context.Context, dto.TaskHistoryFilter) ([]dto.TaskHistory, error)
	Start(int) error
}

type taskService struct {
	service Service
	manager *task.Manager

	task repository.Task
}

// Interface compliance
var _ Task = (*taskService)(nil)

func (s *taskService) GetAll() ([]dto.Task, error) {
	tasks, err := s.manager.Tasks()
	if err != nil {
		return nil, err
	}

	return util.SliceMap(tasks, dto.TaskDTO), nil
}

func (s *taskService) GetHistory(ctx context.Context, filters dto.TaskHistoryFilter) ([]dto.TaskHistory, error) {
	tasks, err := s.task.GetFiltered(ctx, model.TaskFilter(filters))
	if err != nil {
		return nil, err
	}

	return util.SliceMap(tasks, dto.TaskHistoryDTO), nil
}

func (s *taskService) Start(id int) error {
	return s.manager.Run(id)
}
