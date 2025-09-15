package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shareed2k/goth_fiber"
	"go.uber.org/zap"
)

func RoleRoute(c *fiber.Ctx) error {
	session, err := goth_fiber.SessionStore.Get(c)
	if err != nil {
		zap.S().Errorf("failed to get session %v", err)
		return fiber.ErrInternalServerError
	}

	var role any
	if role = session.Get("role"); session == nil {
		return c.Redirect("/", fiber.StatusForbidden)
	}

	// This setup will make it easy to support different roles
	// Think about
	//   - Event organizer
	//   - "Normal" zeus user
	if role == "" {
		return c.Redirect("/", fiber.StatusForbidden)
	}

	return c.Next()
}
