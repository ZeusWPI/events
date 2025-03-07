package service

import (
	"context"
	"fmt"

	"github.com/ZeusWPI/events/internal/api/dto"
	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/pkg/util"
)

// Event represents all business logic regarding events
type Event interface {
	GetByYear(context.Context, dto.Year) ([]dto.Event, error)
	UpdateOrganizers(context.Context, []dto.Event) error
}

type eventService struct {
	service Service

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

	return s.service.withRollback(ctx, func(c context.Context) error {
		for _, event := range events {
			eventDB, found := util.SliceFind(eventsDB, func(e *model.Event) bool { return event.ID == e.ID })
			if !found {
				return fmt.Errorf("Unable to find event %+v", event)
			}

			// Add new organizers
			for _, organizer := range event.Organizers {
				if _, found := util.SliceFind(eventDB.Organizers, func(b model.Board) bool { return organizer.ID == b.ID }); !found {
					if err := s.organizer.Save(c, &model.Organizer{Board: model.Board{ID: organizer.ID}, Event: model.Event{ID: event.ID}}); err != nil {
						return err
					}
				}
			}

			// Remove old organizers
			for _, organizer := range eventDB.Organizers {
				if _, found := util.SliceFind(event.Organizers, func(o dto.Organizer) bool { return organizer.ID == o.ID }); !found {
					if err := s.organizer.Delete(c, model.Organizer{Board: model.Board{ID: organizer.ID}, Event: model.Event{ID: event.ID}}); err != nil {
						return err
					}
				}
			}
		}

		return nil
	})
}
