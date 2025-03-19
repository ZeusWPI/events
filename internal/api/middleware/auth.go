// Package middleware provides various middlewares
package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shareed2k/goth_fiber"
	"go.uber.org/zap"
)

// ProtectedRoute only allows authenticated users to go through
func ProtectedRoute(c *fiber.Ctx) error {
	session, err := goth_fiber.SessionStore.Get(c)
	if err != nil {
		zap.S().Errorf("failed to get session %v", err)
		return fiber.ErrInternalServerError
	}

	if session.Fresh() {
		return c.Redirect("/")
	}

	var userID interface{}
	if userID = session.Get("memberID"); userID == nil {
		return c.SendStatus(fiber.StatusForbidden)
	}

	c.Locals("memberID", userID)

	return c.Next()
}
