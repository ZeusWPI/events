package service

import (
	"context"

	"github.com/ZeusWPI/events/internal/api/dto"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/pkg/util"
)

// Organizer represents all business logic regarding organizers
type Organizer interface {
	GetByYear(context.Context, dto.Year) ([]dto.Organizer, error)
}

type organizerService struct {
	service Service

	board repository.Board
}

// Interface compliance
var _ Organizer = (*organizerService)(nil)

// GetByYear returns all possible organizers for a given year
func (s organizerService) GetByYear(ctx context.Context, year dto.Year) ([]dto.Organizer, error) {
	organizers, err := s.board.GetByYearWithMemberYear(ctx, *year.ToModel())
	if err != nil {
		return nil, err
	}

	return util.SliceMap(organizers, dto.OrganizerDTO), nil
}
