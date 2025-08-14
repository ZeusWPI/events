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

type event struct {
	ID          int        `json:"id"`
	URL         string     `json:"url"` // URL to the website page
	Name        string     `json:"name"`
	Description string     `json:"description"` // Can be empty
	StartTime   time.Time  `json:"start_time"`
	EndTime     *time.Time `json:"end_time,omitempty"` // Might not be present if not applicable
	Location    string     `json:"location"`           // Can be empty
	YearStart   int        `json:"year_start"`
	YearEnd     int        `json:"year_end"`
}

func toEvent(e dto.Event) event {
	return event{
		ID:          e.ID,
		URL:         e.URL,
		Name:        e.Name,
		Description: e.Description,
		StartTime:   e.StartTime,
		EndTime:     e.EndTime,
		Location:    e.Location,
		YearStart:   e.Year.Start,
		YearEnd:     e.Year.End,
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

// get returns all events for the current academic year
//
//	@Summary		Get this years events
//	@Description	Get all planned events for the current academic year.
//	@Tags			event
//	@Produce		json
//	@Success		200	{array}	event
//	@Failure		500
//	@Router			/event [get]
func (r *Event) get(c *fiber.Ctx) error {
	events, err := r.event.GetByLastYear(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(utils.SliceMap(events, toEvent))
}

// getNext returns the next event
//
//	@Summary		Get next event
//	@Description	Get the next event. Returns a 404 if there's no next event planned.
//	@Tags			event
//	@Produce		json
//	@Success		200	{object}	event
//	@Failure		404
//	@Failure		500
//	@Router			/event/next [get]
func (r *Event) getNext(c *fiber.Ctx) error {
	event, err := r.event.GetNext(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(toEvent(event))
}

// getPoster returns the poster for an event
//
//	@Summary		Get event poster
//	@Description	Get the poster for an event. Returns 400 if the event isn't found and 404 if the event doesn't have the requested poster type
//	@Tags			event
//	@Produce		png
//	@Success		200
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Param			id	path	int		true	"event id"
//	@Param			scc	query	boolean	false	"set to true if the scc poster version is desired"
//	@Router			/event/{id} [get]
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
