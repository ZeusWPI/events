package service

import (
	"context"
	"slices"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/server/dto"
	"github.com/ZeusWPI/events/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Event struct {
	service Service

	board        repository.Board
	event        repository.Event
	announcement repository.Announcement
	organizer    repository.Organizer
}

func (s *Service) NewEvent() *Event {
	return &Event{
		service:      *s,
		board:        *s.repo.NewBoard(),
		event:        *s.repo.NewEvent(),
		announcement: *s.repo.NewAnnouncement(),
		organizer:    *s.repo.NewOrganizer(),
	}
}

func (e *Event) GetByYear(ctx context.Context, yearID int) ([]dto.Event, error) {
	eventsDB, err := e.event.GetByYearPopulated(ctx, yearID)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}
	if eventsDB == nil {
		return []dto.Event{}, nil
	}
	events := utils.SliceMap(eventsDB, dto.EventDTO)

	// Add checks

	checks, err := e.service.check.Status(ctx, yearID)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}

	for idx, event := range events {
		if check, ok := checks[event.ID]; ok {
			events[idx].Checks = utils.SliceMap(check, dto.CheckDTO)
		}
	}

	// Add announcements
	announcements, err := e.announcement.GetByEvents(ctx, utils.SliceDereference(eventsDB))
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}

	for _, announcement := range announcements {
		if idx := slices.IndexFunc(events, func(e dto.Event) bool { return e.ID == announcement.EventID }); idx != -1 {
			events[idx].Announcement = dto.AnnouncementDTO(*announcement)
		}
	}

	return events, nil
}

func (e *Event) UpdateOrganizers(ctx context.Context, events []dto.EventOrganizers) error {
	if len(events) == 0 {
		return nil
	}

	eventsDB, err := e.event.GetByIDs(ctx, utils.SliceMap(events, func(e dto.EventOrganizers) int { return e.EventID }))
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	boardsDB, err := e.board.GetByIDs(ctx, utils.SliceFlatten(utils.SliceMap(events, func(e dto.EventOrganizers) []int { return e.Organizers })))
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return e.service.withRollback(ctx, func(ctx context.Context) error {
		for _, event := range events {
			if _, found := utils.SliceFind(eventsDB, func(e *model.Event) bool { return e.ID == event.EventID }); !found {
				zap.S().Debugf("Cant find event %+v", event)
				return fiber.ErrBadRequest
			}

			for _, organizer := range event.Organizers {
				if _, found := utils.SliceFind(boardsDB, func(b *model.Board) bool { return b.ID == organizer }); !found {
					zap.S().Debugf("Cant find organizer %+v", event)
					return fiber.ErrBadRequest
				}
			}

			if err := e.organizer.DeleteByEvent(ctx, event.EventID); err != nil {
				zap.S().Error(err)
				return fiber.ErrInternalServerError
			}

			organizers := make([]model.Organizer, 0, len(event.Organizers))
			for _, organizer := range event.Organizers {
				organizers = append(organizers, model.Organizer{
					EventID: event.EventID,
					BoardID: organizer,
				})
			}

			if err := e.organizer.CreateBatch(ctx, organizers); err != nil {
				zap.S().Error(err)
				return fiber.ErrInternalServerError
			}

		}

		return nil
	})
}
