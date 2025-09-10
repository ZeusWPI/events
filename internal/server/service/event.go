package service

import (
	"context"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/server/dto"
	"github.com/ZeusWPI/events/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Event struct {
	service Service

	announcement repository.Announcement
	board        repository.Board
	check        repository.Check
	event        repository.Event
	mail         repository.Mail
	organizer    repository.Organizer
	poster       repository.Poster
	year         repository.Year
}

func (s *Service) NewEvent() *Event {
	return &Event{
		service:      *s,
		announcement: *s.repo.NewAnnouncement(),
		board:        *s.repo.NewBoard(),
		check:        *s.repo.NewCheck(),
		event:        *s.repo.NewEvent(),
		mail:         *s.repo.NewMail(),
		organizer:    *s.repo.NewOrganizer(),
		poster:       *s.repo.NewPoster(),
		year:         *s.repo.NewYear(),
	}
}

func (e *Event) GetByID(ctx context.Context, eventID int) (dto.Event, error) {
	eventDB, err := e.event.GetByID(ctx, eventID)
	if err != nil {
		zap.S().Error(err)
		return dto.Event{}, fiber.ErrInternalServerError
	}
	if eventDB == nil {
		return dto.Event{}, fiber.ErrBadRequest
	}
	event := dto.EventDTO(eventDB)

	// Add organizers
	organizers, err := e.organizer.GetByEvents(ctx, []model.Event{*eventDB})
	if err != nil {
		zap.S().Error(err)
		return dto.Event{}, fiber.ErrInternalServerError
	}
	event.Organizers = utils.SliceMap(organizers, func(organizer *model.Organizer) dto.Organizer { return dto.OrganizerDTO(&organizer.Board) })

	// Add checks
	checks, err := e.check.GetByEvents(ctx, []model.Event{*eventDB})
	if err != nil {
		zap.S().Error(err)
		return dto.Event{}, fiber.ErrInternalServerError
	}
	event.Checks = utils.SliceMap(checks, dto.CheckDTO)

	// Add posters
	posters, err := e.poster.GetByEvents(ctx, []model.Event{*eventDB})
	if err != nil {
		zap.S().Error(err)
		return dto.Event{}, fiber.ErrInternalServerError
	}
	event.Posters = utils.SliceMap(posters, dto.PosterDTO)

	// Add announcements
	announcements, err := e.announcement.GetByEvents(ctx, []model.Event{*eventDB})
	if err != nil {
		zap.S().Error(err)
		return dto.Event{}, fiber.ErrInternalServerError
	}
	event.Announcements = utils.SliceMap(announcements, dto.AnnouncementDTO)

	// Add mails
	mails, err := e.mail.GetByEvents(ctx, []model.Event{*eventDB})
	if err != nil {
		zap.S().Error(err)
		return dto.Event{}, fiber.ErrInternalServerError
	}
	event.Mails = utils.SliceMap(mails, dto.MailDTO)

	return event, nil
}

func (e *Event) GetNext(ctx context.Context) (dto.Event, error) {
	event, err := e.event.GetNext(ctx)
	if err != nil {
		zap.S().Error(err)
		return dto.Event{}, fiber.ErrInternalServerError
	}
	if event == nil {
		return dto.Event{}, fiber.ErrNotFound
	}

	return e.GetByID(ctx, event.ID)
}

func (e *Event) GetByYear(ctx context.Context, yearID int) ([]dto.Event, error) {
	eventsDB, err := e.event.GetByYear(ctx, yearID)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}
	if eventsDB == nil {
		return []dto.Event{}, nil
	}

	events := make(map[int]dto.Event)
	for _, event := range eventsDB {
		events[event.ID] = dto.EventDTO(event)
	}

	// Add organizers
	organizers, err := e.organizer.GetByEvents(ctx, utils.SliceDereference(eventsDB))
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}

	for _, organizer := range organizers {
		event, ok := events[organizer.EventID]
		if !ok {
			// Should not be possible
			zap.S().Error("Unknown error occurred\nOrganizer queried that has an unknown event %+v | %+v", events, utils.SliceDereference(organizers))
			return nil, fiber.ErrInternalServerError
		}
		event.Organizers = append(event.Organizers, dto.OrganizerDTO(&organizer.Board))
		events[organizer.EventID] = event
	}

	// Add checks
	checks, err := e.check.GetByEvents(ctx, utils.SliceDereference(eventsDB))
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}

	for _, check := range checks {
		event, ok := events[check.EventID]
		if !ok {
			// Should not be possible
			zap.S().Error("Unknown error occured\nCheck queried that has an unknown event %+v | %+v", events, utils.SliceDereference(checks))
			return nil, fiber.ErrInternalServerError
		}
		event.Checks = append(event.Checks, dto.CheckDTO(check))
		events[check.EventID] = event
	}

	return utils.MapValues(events), nil
}

func (e *Event) GetByLastYear(ctx context.Context) ([]dto.Event, error) {
	year, err := e.year.GetLast(ctx)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}
	if year == nil {
		return nil, nil
	}

	return e.GetByYear(ctx, year.ID)
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
				return fiber.ErrBadRequest
			}

			for _, organizer := range event.Organizers {
				if _, found := utils.SliceFind(boardsDB, func(b *model.Board) bool { return b.ID == organizer }); !found {
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
