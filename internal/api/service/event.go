package service

import (
	"context"
	"fmt"

	"github.com/ZeusWPI/events/internal/api/dto"
	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/website"
	"github.com/ZeusWPI/events/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Event struct {
	service Service

	board     repository.Board
	event     repository.Event
	organizer repository.Organizer
}

func newEvent(service Service) *Event {
	return &Event{
		service:   service,
		board:     *service.repo.NewBoard(),
		event:     *service.repo.NewEvent(),
		organizer: *service.repo.NewOrganizer(),
	}
}

func (e *Event) GetByYear(ctx context.Context, id int) ([]dto.Event, error) {
	events, err := e.event.GetByYearPopulated(ctx, id)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}
	if events == nil {
		return []dto.Event{}, nil
	}

	return utils.SliceMap(events, dto.EventDTO), nil
}

func (e *Event) UpdateOrganizers(ctx context.Context, events []dto.Event) error {
	if len(events) == 0 {
		return nil
	}

	eventsDB, err := e.event.GetByYearPopulated(ctx, events[0].Year.ID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}
	if eventsDB == nil {
		eventsDB = []*model.Event{}
	}

	boardsDB, err := e.board.GetByYearPopulated(ctx, events[0].Year.ID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}
	if boardsDB == nil {
		boardsDB = []*model.Board{}
	}

	return e.service.withRollback(ctx, func(c context.Context) error {
		for _, event := range events {
			eventDB, found := utils.SliceFind(eventsDB, func(e *model.Event) bool { return event.ID == e.ID })
			if !found {
				return fmt.Errorf("find event %+v", event)
			}

			boards := make([]model.Board, 0, len(event.Organizers))
			for _, organizer := range event.Organizers {
				board, found := utils.SliceFind(boardsDB, func(b *model.Board) bool { return b.Member.ID == organizer.ID })
				if !found {
					return fmt.Errorf("find given organizer for event %+v | %+v", organizer, event)
				}
				boards = append(boards, *board)
			}

			for _, board := range boards {
				if _, found := utils.SliceFind(eventDB.Organizers, func(b model.Board) bool { return board.ID == b.ID }); !found {
					if err := e.organizer.Create(c, board.ID, event.ID); err != nil {
						return err
					}
				}
			}

			// Remove old organizers
			for _, organizer := range eventDB.Organizers {
				if _, found := utils.SliceFind(boards, func(b model.Board) bool { return organizer.ID == b.ID }); !found {
					if err := e.organizer.DeleteByBoardEvent(c, organizer.ID, event.ID); err != nil {
						return err
					}
				}
			}
		}

		return nil
	})
}

func (e *Event) Sync() error {
	if err := e.service.manager.RunByName(website.YearTask); err != nil {
		return err
	}

	return nil
}
