package service

import (
	"context"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/server/dto"
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

func (s *Service) NewOrganizer() *Organizer {
	return &Organizer{
		service: *s,
		board:   *s.repo.NewBoard(),
		member:  *s.repo.NewMember(),
		year:    *s.repo.NewYear(),
	}
}

func (o *Organizer) GetByMember(ctx context.Context, memberID int) (dto.Organizer, error) {
	member, err := o.member.GetByID(ctx, memberID)
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

	organizers = utils.SliceFilter(organizers, func(b *model.Board) bool { return b.IsOrganizer })

	return utils.SliceMap(organizers, dto.OrganizerDTO), nil
}

func (o *Organizer) GetByZauth(ctx context.Context, zauth zauth.User) (dto.Organizer, error) {
	member, err := o.member.GetByName(ctx, zauth.FullName)
	if err != nil {
		zap.S().Error(err)
		return dto.Organizer{}, err
	}
	if member == nil {
		// Member not in DB yet
		// Probably means he's not a board member
		// If he turns out to be a board member it will be corrected by a bestuur update task
		// Let's add the user as a non organizer member
		if err := o.service.withRollback(ctx, func(ctx context.Context) error {
			member = &model.Member{
				Name:     zauth.FullName,
				Username: zauth.Username,
				ZauthID:  zauth.ID,
			}

			if err = o.member.Create(ctx, member); err != nil {
				zap.S().Error(err)
				return fiber.ErrInternalServerError
			}

			year, err := o.year.GetLast(ctx)
			if err != nil {
				zap.S().Error(err)
				return fiber.ErrInternalServerError
			}

			board := model.Board{
				MemberID:    member.ID,
				YearID:      year.ID,
				Role:        "Niet bestuur",
				IsOrganizer: false,
			}

			if err := o.board.Create(ctx, &board); err != nil {
				zap.S().Error(err)
				return fiber.ErrInternalServerError
			}

			return nil
		}); err != nil {
			return dto.Organizer{}, err
		}
	}

	if member.ZauthID == 0 {
		member.ZauthID = zauth.ID
		member.Username = zauth.Username

		if err = o.member.Update(ctx, *member); err != nil {
			zap.S().Error(err)
			return dto.Organizer{}, err
		}
	}

	return dto.Organizer{
		ID:   member.ID,
		Name: member.Name,
	}, nil
}
