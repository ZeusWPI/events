package api

import (
	"github.com/ZeusWPI/events/internal/api/service"
	"github.com/ZeusWPI/events/pkg/config"
	"github.com/ZeusWPI/events/pkg/util"
	"github.com/ZeusWPI/events/pkg/zauth"
	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
	"github.com/shareed2k/goth_fiber"
	"go.uber.org/zap"
)

// Auth contains all api routes related to authentication
type Auth struct {
	router fiber.Router

	organizer service.Organizer

	redirectURL string
}

// NewAuth creates a new authentication router
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
		organizer:   service.NewOrganizer(),
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

	zauthAdmin, err1 := util.MapGetKeyAsType[bool]("admin", user.RawData)
	zauthID, err2 := util.MapGetKeyAsType[int]("id", user.RawData)
	zauthName, err3 := util.MapGetKeyAsType[string]("fullName", user.RawData)
	zauthUsername, err4 := util.MapGetKeyAsType[string]("username", user.RawData)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		zap.S().Error(err1, err2, err3, err4)
		return fiber.ErrInternalServerError
	}

	session, err := goth_fiber.SessionStore.Get(c)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	if zauthAdmin {
		zauth := zauth.User{
			ID:       zauthID,
			Admin:    zauthAdmin,
			FullName: zauthName,
			Username: zauthUsername,
		}

		dbUser, err := r.organizer.GetByZauth(c.Context(), zauth)
		if err != nil {
			zap.S().Error(err)
			return fiber.ErrInternalServerError
		}

		if dbUser.ID != 0 {
			session.Set("memberID", dbUser.ID)
		}
	}

	if err = session.Save(); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return c.Redirect(r.redirectURL)
}

func (r *Auth) logoutHandler(c *fiber.Ctx) error {
	if err := goth_fiber.Logout(c); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	session, err := goth_fiber.SessionStore.Get(c)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	if err := session.Destroy(); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return c.SendStatus(fiber.StatusOK)
}
