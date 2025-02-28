package api

import (
	"github.com/ZeusWPI/events/internal/api/dto"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/pkg/util"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// YearRouter contains all api routes related to year
type YearRouter struct {
	router fiber.Router

	year repository.Year
}

// NewYearRouter constructs a new year router
func NewYearRouter(repo repository.Repository, router fiber.Router) *YearRouter {
	api := &YearRouter{
		router: router.Group("/year"),
		year:   repo.NewYear(),
	}
	api.createRoutes()

	return api
}

func (r *YearRouter) createRoutes() {
	r.router.Get("/", r.getAll)
}

func (r *YearRouter) getAll(c *fiber.Ctx) error {
	years, err := r.year.GetAll(c.Context())
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(util.SliceMap(years, dto.YearDTO))
}
