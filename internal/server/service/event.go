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

	announcement repository.Announcement
	board        repository.Board
	event        repository.Event
	mail         repository.Mail
	organizer    repository.Organizer
	year         repository.Year
}

func (s *Service) NewEvent() *Event {
	return &Event{
		service:      *s,
		announcement: *s.repo.NewAnnouncement(),
		board:        *s.repo.NewBoard(),
		event:        *s.repo.NewEvent(),
		mail:         *s.repo.NewMail(),
		organizer:    *s.repo.NewOrganizer(),
		year:         *s.repo.NewYear(),
	}
}

func (e *Event) GetByID(ctx context.Context, eventID int) (dto.Event, error) {
	eventDB, err := e.event.GetByIDPopulated(ctx, eventID)
	if err != nil {
		zap.S().Error(err)
		return dto.Event{}, fiber.ErrInternalServerError
	}
	if eventDB == nil {
		return dto.Event{}, fiber.ErrBadRequest
	}
	event := dto.EventDTO(eventDB)

	// Add checks

	checks, err := e.service.check.Status(ctx, eventDB.YearID)
	if err != nil {
		zap.S().Error(err)
		return dto.Event{}, fiber.ErrInternalServerError
	}

	if check, ok := checks[event.ID]; ok {
		event.Checks = utils.SliceMap(check, dto.CheckDTO)
	}

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
	event, err := e.event.GetNextWithYear(ctx)
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
		for _, eventID := range announcement.EventIDs {
			if idx := slices.IndexFunc(events, func(e dto.Event) bool { return e.ID == eventID }); idx != -1 {
				events[idx].Announcements = append(events[idx].Announcements, dto.AnnouncementDTO(announcement))
			}
		}
	}

	// Add mails
	mails, err := e.mail.GetByEvents(ctx, utils.SliceDereference(eventsDB))
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}

	for _, mail := range mails {
		for _, eventID := range mail.EventIDs {
			if idx := slices.IndexFunc(events, func(e dto.Event) bool { return e.ID == eventID }); idx != -1 {
				events[idx].Mails = append(events[idx].Mails, dto.MailDTO(mail))
			}
		}
	}

	return events, nil
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
