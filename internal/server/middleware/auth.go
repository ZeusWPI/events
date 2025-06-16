package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shareed2k/goth_fiber"
	"go.uber.org/zap"
)

func ProtectedRoute(c *fiber.Ctx) error {
	session, err := goth_fiber.SessionStore.Get(c)
	if err != nil {
		zap.S().Errorf("failed to get session %v", err)
		return fiber.ErrInternalServerError
	}

	if session.Fresh() {
		return c.Redirect("/", fiber.StatusUnauthorized)
	}

	var userID interface{}
	if userID = session.Get("memberID"); userID == nil {
		return c.Redirect("/", fiber.StatusForbidden)
	}

	c.Locals("memberID", userID)

	return c.Next()
}
