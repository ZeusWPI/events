package service

import (
	"context"

	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/server/dto"
	"github.com/ZeusWPI/events/pkg/utils"
)

type Year struct {
	service Service

	year repository.Year
}

func (s *Service) NewYear() *Year {
	return &Year{
		service: *s,
		year:    *s.repo.NewYear(),
	}
}

func (y *Year) GetAll(ctx context.Context) ([]dto.Year, error) {
	years, err := y.year.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	if years == nil {
		return []dto.Year{}, nil
	}

	return utils.SliceMap(years, dto.YearDTO), nil
}
