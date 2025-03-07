package service

import (
	"context"

	"github.com/ZeusWPI/events/internal/api/dto"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/pkg/util"
)

// Year represents all business logic regarding years
type Year interface {
	GetAll(context.Context) ([]dto.Year, error)
}

type yearService struct {
	service Service

	year repository.Year
}

// Interface compliance
var _ Year = (*yearService)(nil)

// GetAll returns all years
func (s *yearService) GetAll(ctx context.Context) ([]dto.Year, error) {
	years, err := s.year.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return util.SliceMap(years, dto.YearDTO), nil
}
