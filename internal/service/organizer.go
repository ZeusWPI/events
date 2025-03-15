package service

import (
	"context"

	"github.com/ZeusWPI/events/internal/api/dto"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/pkg/util"
	"github.com/ZeusWPI/events/pkg/zauth"
)

// Organizer represents all business logic regarding organizers
type Organizer interface {
	GetByID(context.Context, int) (dto.Organizer, error)
	GetByYear(context.Context, dto.Year) ([]dto.Organizer, error)
	GetByZauth(context.Context, zauth.User) (dto.Organizer, error)
}

type organizerService struct {
	service Service

	board  repository.Board
	member repository.Member
	year   repository.Year
}

// Interface compliance
var _ Organizer = (*organizerService)(nil)

// GetByID returns an organizer given an id.
// It's associated role is the one for the current year or empty
func (s *organizerService) GetByID(ctx context.Context, id int) (dto.Organizer, error) {
	member, err := s.member.GetByID(ctx, id)
	if err != nil {
		return dto.Organizer{}, err
	}

	year, err := s.year.GetLatest(ctx)
	if err != nil {
		return dto.Organizer{}, err
	}

	board, err := s.board.GetByMemberYear(ctx, *member, *year)
	if err != nil {
		return dto.Organizer{}, err
	}

	if board.ID == 0 {
		// No board entry found, manually populate
		board.Member = *member
	}

	organizer := dto.OrganizerDTO(board)

	return organizer, nil
}

// GetByYear returns all possible organizers for a given year
func (s *organizerService) GetByYear(ctx context.Context, year dto.Year) ([]dto.Organizer, error) {
	organizers, err := s.board.GetByYearWithMemberYear(ctx, *year.ToModel())
	if err != nil {
		return nil, err
	}

	return util.SliceMap(organizers, dto.OrganizerDTO), nil
}

// GetByZauth returns the organizer associated with a zauth user
// If no organizer is found the returned struct will have an ID of 0
func (s *organizerService) GetByZauth(ctx context.Context, zauth zauth.User) (dto.Organizer, error) {
	member, err := s.member.GetByName(ctx, zauth.FullName)
	if err != nil {
		return dto.Organizer{}, err
	}

	if member.ZauthID == 0 {
		// First time querying for this user
		member.ZauthID = zauth.ID
		member.Username = zauth.Username

		if err = s.member.Save(ctx, member); err != nil {
			return dto.Organizer{}, err
		}
	}

	return dto.Organizer{
		ID:   member.ID,
		Name: member.Name,
	}, nil
}
