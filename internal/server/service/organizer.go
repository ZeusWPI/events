package service

import (
	"context"
	"slices"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/server/dto"
	"github.com/ZeusWPI/events/pkg/config"
	"github.com/ZeusWPI/events/pkg/utils"
	"github.com/ZeusWPI/events/pkg/zauth"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Organizer struct {
	service     Service
	development bool

	board  repository.Board
	member repository.Member
	year   repository.Year
}

func (s *Service) NewOrganizer() *Organizer {
	return &Organizer{
		service:     *s,
		development: config.IsDev(),
		board:       *s.repo.NewBoard(),
		member:      *s.repo.NewMember(),
		year:        *s.repo.NewYear(),
	}
}

func (o *Organizer) GetByMember(ctx context.Context, memberID int) (dto.Organizer, error) {
	member, err := o.member.GetByID(ctx, memberID)
	if err != nil {
		zap.S().Error(err)
		return dto.Organizer{}, fiber.ErrInternalServerError
	}
	if member == nil {
		return dto.Organizer{}, fiber.ErrNotFound
	}

	year, err := o.year.GetLast(ctx)
	if err != nil {
		zap.S().Error(err)
		return dto.Organizer{}, fiber.ErrInternalServerError
	}

	board, err := o.board.GetByMemberYear(ctx, member.ID, year.ID)
	if err != nil {
		zap.S().Error(err)
		return dto.Organizer{}, fiber.ErrInternalServerError
	}
	if board == nil {
		return dto.Organizer{}, fiber.ErrNotFound
	}
	board.Member = *member

	return dto.OrganizerDTO(board), nil
}

func (o *Organizer) GetByYear(ctx context.Context, yearID int) ([]dto.Organizer, error) {
	organizers, err := o.board.GetByYear(ctx, yearID)
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
		// Probably means the user was never a board member
		// If he turns out to be a board member it will be corrected by a bestuur update task
		// Let's add the user
		member = &model.Member{
			Name:     zauth.FullName,
			Username: zauth.Username,
			ZauthID:  zauth.ID,
		}

		if err = o.member.Create(ctx, member); err != nil {
			zap.S().Error(err)
			return dto.Organizer{}, fiber.ErrInternalServerError
		}
	}
	year, err := o.year.GetLast(ctx)
	if err != nil {
		zap.S().Error(err)
		return dto.Organizer{}, fiber.ErrInternalServerError
	}

	board, err := o.board.GetByMemberYear(ctx, member.ID, year.ID)
	if err != nil {
		zap.S().Error(err)
		return dto.Organizer{}, fiber.ErrInternalServerError
	}

	if o.development && board == nil {
		// Development environment
		// Add user to the board to give permission
		// to edit everything
		board = &model.Board{
			MemberID:    member.ID,
			YearID:      year.ID,
			Role:        "Developer",
			IsOrganizer: true,
		}

		if err := o.board.Create(ctx, board); err != nil {
			zap.S().Error(err)
			return dto.Organizer{}, fiber.ErrInternalServerError
		}
	}

	if slices.Contains(zauth.Roles, "events_admin") && board == nil {
		// User is an events admin
		// Add the user to the board
		board = &model.Board{
			MemberID:    member.ID,
			YearID:      year.ID,
			Role:        "Events Admin",
			IsOrganizer: false,
		}

		if err := o.board.Create(ctx, board); err != nil {
			zap.S().Error(err)
			return dto.Organizer{}, fiber.ErrInternalServerError
		}
	}

	if member.ZauthID == 0 {
		// Member already exists but it's the first time the user logs in
		// Probably means the user is / was a board member and was added by the bestuur update task
		// Let's link zauth id with their name
		member.ZauthID = zauth.ID
		member.Username = zauth.Username

		if err = o.member.Update(ctx, *member); err != nil {
			zap.S().Error(err)
			return dto.Organizer{}, err
		}
	}

	role := ""
	if board != nil {
		role = board.Role
	}

	return dto.Organizer{
		ID:   member.ID,
		Name: member.Name,
		Role: role,
	}, nil
}
