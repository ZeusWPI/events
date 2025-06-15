package service

import (
	"context"

	"github.com/ZeusWPI/events/internal/api/dto"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/pkg/utils"
	"github.com/ZeusWPI/events/pkg/zauth"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Organizer struct {
	service Service

	board  repository.Board
	member repository.Member
	year   repository.Year
}

func newOrganizer(service Service) *Organizer {
	return &Organizer{
		service: service,
		board:   *service.repo.NewBoard(),
		member:  *service.repo.NewMember(),
		year:    *service.repo.NewYear(),
	}
}

func (o *Organizer) GetByID(ctx context.Context, id int) (dto.Organizer, error) {
	member, err := o.member.GetByID(ctx, id)
	if err != nil {
		zap.S().Error(err)
		return dto.Organizer{}, fiber.ErrInternalServerError
	}
	if member == nil {
		return dto.Organizer{}, fiber.ErrBadRequest
	}

	year, err := o.year.GetLast(ctx)
	if err != nil {
		zap.S().Error(err)
		return dto.Organizer{}, fiber.ErrInternalServerError
	}

	board, err := o.board.GetByMemberYear(ctx, *member, *year)
	if err != nil {
		zap.S().Error(err)
		return dto.Organizer{}, fiber.ErrInternalServerError
	}
	if board == nil {
		return dto.Organizer{}, fiber.ErrBadRequest
	}

	return dto.OrganizerDTO(board), nil
}

func (o *Organizer) GetByYear(ctx context.Context, yearID int) ([]dto.Organizer, error) {
	organizers, err := o.board.GetByYearPopulated(ctx, yearID)
	if err != nil {
		return nil, err
	}
	if organizers == nil {
		return []dto.Organizer{}, nil
	}

	return utils.SliceMap(organizers, dto.OrganizerDTO), nil
}

func (o *Organizer) GetByZauth(ctx context.Context, zauth zauth.User) (dto.Organizer, error) {
	member, err := o.member.GetByName(ctx, zauth.FullName)
	if err != nil {
		zap.S().Error(err)
		return dto.Organizer{}, err
	}
	if member == nil {
		return dto.Organizer{}, fiber.ErrBadRequest
	}

	if member.ZauthID == 0 {
		// First time querying for this user
		member.ZauthID = zauth.ID
		member.Username = zauth.Username

		if err := o.member.Create(ctx, member); err != nil {
			zap.S().Error(err)
			return dto.Organizer{}, err
		}
	}

	return dto.Organizer{
		ID:   member.ID,
		Name: member.Name,
	}, nil
}
