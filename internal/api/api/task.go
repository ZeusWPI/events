package api

import (
	"strconv"

	"github.com/ZeusWPI/events/internal/api/dto"
	"github.com/ZeusWPI/events/internal/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// TaskRouter contains all api routes regarding tasks
type TaskRouter struct {
	router fiber.Router

	task service.Task
}

// NewTaskRouter creates a new task router
func NewTaskRouter(service service.Service, router fiber.Router) *TaskRouter {
	api := &TaskRouter{
		router: router.Group("/task"),
		task:   service.NewTask(),
	}

	api.createRoutes()

	return api
}

func (r *TaskRouter) createRoutes() {
	r.router.Get("/history", r.getHistory)
	r.router.Post("/:id", r.start)
	r.router.Get("/", r.getAll)
}

func (r *TaskRouter) getAll(c *fiber.Ctx) error {
	tasks, err := r.task.GetAll()
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(tasks)
}

func (r *TaskRouter) getHistory(c *fiber.Ctx) error {
	name := c.Query("name")
	onlyErrored := c.Query("only_errored", "false") == "true"

	var recurring *bool
	recurringRaw := c.Query("recurring")
	if recurringRaw != "" {
		recurringValue := recurringRaw == "true"
		recurring = &recurringValue
	}

	page, err1 := strconv.Atoi(c.Query("page", "1"))
	limit, err2 := strconv.Atoi(c.Query("limit", "10"))
	if err1 != nil || err2 != nil || page < 0 || limit < 1 {
		return fiber.ErrBadRequest
	}

	tasks, err := r.task.GetHistory(c.Context(), dto.TaskHistoryFilter{
		Name:        name,
		OnlyErrored: onlyErrored,
		Recurring:   recurring,
		Page:        page,
		Limit:       limit,
	})
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(tasks)
}

func (r *TaskRouter) start(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	if err := r.task.Start(id); err != nil {
		return fiber.ErrInternalServerError
	}

	return c.SendStatus(fiber.StatusAccepted)
}
