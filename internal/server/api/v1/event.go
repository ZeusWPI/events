// Package v1 contains the first version of the public API
package v1

import (
	"time"

	"github.com/ZeusWPI/events/internal/server/dto"
	"github.com/ZeusWPI/events/internal/server/service"
	"github.com/ZeusWPI/events/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

const mimePNG = "image/png"

type Event struct {
	router fiber.Router

	event  service.Event
	poster service.Poster
}

type eventAPI struct {
	ID          int        `json:"id"`
	URL         string     `json:"url"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     *time.Time `json:"end_time,omitempty"` // Pointer to support omitempty
	Location    string     `json:"location"`
	YearStart   int        `json:"year_start"`
	YearEnd     int        `json:"year_end"`
}

func toEventAPI(event dto.Event) eventAPI {
	return eventAPI{
		ID:          event.ID,
		URL:         event.URL,
		Name:        event.Name,
		Description: event.Description,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
		Location:    event.Location,
		YearStart:   event.Year.Start,
		YearEnd:     event.Year.End,
	}
}

func NewEvent(router fiber.Router, service *service.Service) *Event {
	api := &Event{
		router: router.Group("/event"),
		event:  *service.NewEvent(),
		poster: *service.NewPoster(),
	}

	api.createRoutes()

	return api
}

func (r *Event) createRoutes() {
	r.router.Get("/next", r.getNext)
	r.router.Get("/:id", r.getPoster)
	r.router.Get("/", r.get)
}

func (r *Event) get(c *fiber.Ctx) error {
	events, err := r.event.GetByLastYear(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(utils.SliceMap(events, toEventAPI))
}

func (r *Event) getNext(c *fiber.Ctx) error {
	event, err := r.event.GetNext(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(toEventAPI(event))
}

func (r *Event) getPoster(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	scc := c.QueryBool("scc", false)

	event, err := r.event.GetByID(c.Context(), id)
	if err != nil {
		return err
	}
	if len(event.Posters) == 0 {
		return fiber.ErrNotFound
	}

	poster, ok := utils.SliceFind(event.Posters, func(p dto.Poster) bool { return p.SCC == scc })
	if !ok {
		return fiber.ErrNotFound
	}

	file, err := r.poster.GetFile(c.Context(), poster.ID)
	if err != nil {
		return err
	}

	c.Set("Content-Type", mimePNG)
	return c.Send(file)
}
