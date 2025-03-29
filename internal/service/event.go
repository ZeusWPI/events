package service

import (
	"context"
	"fmt"

	"github.com/ZeusWPI/events/internal/api/dto"
	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/pkg/task"
	"github.com/ZeusWPI/events/internal/pkg/website"
	"github.com/ZeusWPI/events/pkg/util"
)

// Event represents all business logic regarding events
type Event interface {
	GetByYear(context.Context, dto.Year) ([]dto.Event, error)
	UpdateOrganizers(context.Context, []dto.Event) error
	Sync() error
}

type eventService struct {
	service Service
	manager *task.Manager

	board     repository.Board
	event     repository.Event
	organizer repository.Organizer
}

// Interface compliance
var _ Event = (*eventService)(nil)

// GetByYear returns all events occuring in a given year
func (s *eventService) GetByYear(ctx context.Context, year dto.Year) ([]dto.Event, error) {
	events, err := s.event.GetByYearWithAll(ctx, *year.ToModel())
	if err != nil {
		return nil, err
	}

	return util.SliceMap(events, dto.EventDTO), nil
}

// UpdateOrganizers updates the organizers for each given event
func (s *eventService) UpdateOrganizers(ctx context.Context, events []dto.Event) error {
	if len(events) == 0 {
		return nil
	}

	eventsDB, err := s.event.GetByYearWithAll(ctx, *events[0].Year.ToModel())
	if err != nil {
		return err
	}

	boardsDB, err := s.board.GetByYearWithMemberYear(ctx, *events[0].Year.ToModel())
	if err != nil {
		return err
	}

	return s.service.withRollback(ctx, func(c context.Context) error {
		for _, event := range events {
			// Find existing event
			eventDB, found := util.SliceFind(eventsDB, func(e *model.Event) bool { return event.ID == e.ID })
			if !found {
				return fmt.Errorf("unable to find event %+v", event)
			}

			// Find corresponding board members
			boards := make([]model.Board, 0, len(event.Organizers))
			for _, organizer := range event.Organizers {
				board, found := util.SliceFind(boardsDB, func(b *model.Board) bool { return b.Member.ID == organizer.ID })
				if !found {
					return fmt.Errorf("unable to find given organizer for event %+v | %+v", organizer, event)
				}
				boards = append(boards, *board)
			}

			// Add new organizers
			for _, board := range boards {
				if _, found := util.SliceFind(eventDB.Organizers, func(b model.Board) bool { return board.ID == b.ID }); !found {
					if err := s.organizer.Save(c, &model.Organizer{Board: board, Event: model.Event{ID: event.ID}}); err != nil {
						return err
					}
				}
			}

			// Remove old organizers
			for _, organizer := range eventDB.Organizers {
				if _, found := util.SliceFind(boards, func(b model.Board) bool { return organizer.ID == b.ID }); !found {
					if err := s.organizer.Delete(c, model.Organizer{Board: organizer, Event: model.Event{ID: event.ID}}); err != nil {
						return err
					}
				}
			}
		}

		return nil
	})
}

func (s *eventService) Sync() error {
	if err := s.manager.RunByName(website.YearTask); err != nil {
		return err
	}

	return nil
}
