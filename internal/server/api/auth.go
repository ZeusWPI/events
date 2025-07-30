package api

import (
	"fmt"

	"github.com/ZeusWPI/events/internal/server/service"
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
	development bool
}

func NewAuth(service service.Service, router fiber.Router) *Auth {
	goth.UseProviders(
		zauth.NewProvider(
			config.GetString("auth.client"),
			config.GetString("auth.secret"),
			config.GetString("auth.callback_url"),
		),
	)

	api := &Auth{
		router:      router.Group("/auth"),
		organizer:   *service.NewOrganizer(),
		redirectURL: config.GetDefaultString("auth.redirect_url", "/"),
		development: config.GetDefaultString("app.env", "development") == "development",
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

	zauth, err := utils.MapGetKeyAsType[zauth.User]("user", user.RawData)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrBadGateway
	}

	if !r.development {
		// Only restrict application access in non dev environment
		if !utils.SliceContainsAny(zauth.Roles, []string{"bestuur", "events_admin"}) {
			return fiber.ErrForbidden
		}
	}

	dbUser, err := r.organizer.GetByZauth(c.Context(), zauth)
	if err != nil {
		return err
	}

	if err = storeInSession(c, "memberID", dbUser.ID); err != nil {
		zap.S().Errorf("Failed to store member id in session %v", err)
		return fiber.ErrInternalServerError
	}

	zap.S().Debug(dbUser)
	zap.S().Debug("Auth done")

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

func storeInSession(ctx *fiber.Ctx, key string, value any) error {
	session, err := goth_fiber.SessionStore.Get(ctx)
	if err != nil {
		return fmt.Errorf("get session %w", err)
	}

	session.Set(key, value)

	return session.Save()
}
