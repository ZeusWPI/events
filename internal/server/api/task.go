package api

import (
	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/server/dto"
	"github.com/ZeusWPI/events/internal/server/service"
	"github.com/gofiber/fiber/v2"
)

type Task struct {
	router fiber.Router

	task service.Task
}

func NewTask(router fiber.Router, service *service.Service) *Task {
	api := &Task{
		router: router.Group("/task"),
		task:   *service.NewTask(),
	}

	api.createRoutes()

	return api
}

func (r *Task) createRoutes() {
	r.router.Get("/history", r.getHistory)
	r.router.Post("/resolve/:id", r.resolve)
	r.router.Post("/start/:id", r.start)
	r.router.Get("/", r.getAll)
}

func (r *Task) getAll(c *fiber.Ctx) error {
	tasks, err := r.task.GetAll()
	if err != nil {
		return err
	}

	return c.JSON(tasks)
}

func (r *Task) getHistory(c *fiber.Ctx) error {
	uid := c.Query("uid")
	resultStr := c.Query("result")

	var result *model.TaskResult
	switch resultStr {
	case string(model.Success), string(model.Resolved), string(model.Failed):
		resultTmp := model.TaskResult(resultStr)
		result = &resultTmp
	}

	limit := c.QueryInt("limit", 10)
	page := c.QueryInt("page", 0)
	if limit < 1 || page < 0 {
		return fiber.ErrBadRequest
	}

	tasks, err := r.task.GetHistory(c.Context(), dto.TaskFilter{
		TaskUID: uid,
		Result:  result,
		Limit:   limit,
		Offset:  page * limit,
	})
	if err != nil {
		return err
	}

	return c.JSON(tasks)
}

func (r *Task) resolve(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	if err := r.task.Resolve(c.Context(), id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func (r *Task) start(c *fiber.Ctx) error {
	uid := c.Params("uid")
	if err := r.task.Start(uid); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusAccepted)
}
