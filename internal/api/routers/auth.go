package api

import (
	"fmt"

	"github.com/ZeusWPI/events/internal/api/service"
	"github.com/ZeusWPI/events/pkg/config"
	"github.com/ZeusWPI/events/pkg/utils"
	"github.com/ZeusWPI/events/pkg/zauth"
	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
	"github.com/shareed2k/goth_fiber"
	"go.uber.org/zap"
)

type Auth struct {
	router    fiber.Router
	organizer service.Organizer

	redirectURL string
}

func NewAuth(service service.Service, router fiber.Router) *Auth {
	goth.UseProviders(
		zauth.New(
			config.GetString("auth.client"),
			config.GetString("auth.secret"),
			config.GetString("auth.callback_url"),
		),
	)

	api := &Auth{
		router:      router.Group("/auth"),
		organizer:   *service.NewOrganizer(),
		redirectURL: config.GetDefaultString("auth.redirect_url", "/"),
	}
	api.createRoutes()

	return api
}

func (r *Auth) createRoutes() {
	r.router.Get("/login/:provider", goth_fiber.BeginAuthHandler)
	r.router.Get("/callback/:provider", r.loginCallback)
	r.router.Post("/logout", r.logoutHandler)
}

func (r *Auth) loginCallback(c *fiber.Ctx) error {
	user, err := goth_fiber.CompleteUserAuth(c)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	zauthAdmin, err1 := utils.MapGetKeyAsType[bool]("admin", user.RawData)
	zauthID, err2 := utils.MapGetKeyAsType[int]("id", user.RawData)
	zauthName, err3 := utils.MapGetKeyAsType[string]("fullName", user.RawData)
	zauthUsername, err4 := utils.MapGetKeyAsType[string]("username", user.RawData)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		zap.S().Error(err1, err2, err3, err4)
		return fiber.ErrBadGateway
	}

	if !zauthAdmin {
		return fiber.ErrForbidden
	}

	zauth := zauth.User{
		ID:       zauthID,
		Admin:    zauthAdmin,
		FullName: zauthName,
		Username: zauthUsername,
	}

	dbUser, err := r.organizer.GetByZauth(c.Context(), zauth)
	if err != nil {
		// If err == fiber.ErrBadRequest then the user is admin but not a board member
		return err
	}

	if err = storeInSession(c, "memberID", dbUser.ID); err != nil {
		zap.S().Errorf("Failed to store member id in session %v", err)
		return fiber.ErrInternalServerError
	}

	return c.Redirect(r.redirectURL)
}

func (r *Auth) logoutHandler(c *fiber.Ctx) error {
	if err := goth_fiber.Logout(c); err != nil {
		zap.S().Errorf("Failed to logout %v", err)
	}

	session, err := goth_fiber.SessionStore.Get(c)
	if err != nil {
		zap.S().Errorf("Failed to get session %v", err)
		return fiber.ErrInternalServerError
	}
	if err := session.Destroy(); err != nil {
		zap.S().Errorf("Failed to destroy %v", err)
		return fiber.ErrInternalServerError
	}

	return c.SendStatus(fiber.StatusOK)
}

func storeInSession(ctx *fiber.Ctx, key string, value interface{}) error {
	session, err := goth_fiber.SessionStore.Get(ctx)
	if err != nil {
		return fmt.Errorf("get session %w", err)
	}

	session.Set(key, value)

	return session.Save()
}
