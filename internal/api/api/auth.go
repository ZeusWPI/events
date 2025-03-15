package api

import (
	"github.com/ZeusWPI/events/internal/service"
	"github.com/ZeusWPI/events/pkg/config"
	"github.com/ZeusWPI/events/pkg/util"
	"github.com/ZeusWPI/events/pkg/zauth"
	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
	"github.com/shareed2k/goth_fiber"
	"go.uber.org/zap"
)

// AuthRouter contains all api routes related to authentication
type AuthRouter struct {
	router fiber.Router

	organizer service.Organizer

	redirectURL string
}

// NewAuthRouter creates a new authentication router
func NewAuthRouter(service service.Service, router fiber.Router) *AuthRouter {
	goth.UseProviders(
		zauth.New(
			config.GetString("auth.client"),
			config.GetString("auth.secret"),
			config.GetString("auth.callback_url"),
		),
	)

	api := &AuthRouter{
		router:      router.Group("/auth"),
		organizer:   service.NewOrganizer(),
		redirectURL: config.GetDefaultString("auth.redirect_url", "/"),
	}
	api.createRoutes()

	return api
}

func (r *AuthRouter) createRoutes() {
	r.router.Get("/login/:provider", goth_fiber.BeginAuthHandler)
	r.router.Get("/callback/:provider", r.loginCallback)
	r.router.Post("/logout", r.logoutHandler)
}

func (r *AuthRouter) loginCallback(c *fiber.Ctx) error {
	user, err := goth_fiber.CompleteUserAuth(c)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	zauthAdmin, err := util.MapGetKeyAsType[bool]("admin", user.RawData)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}
	if !zauthAdmin {
		return fiber.ErrUnauthorized
	}

	zauthID, err1 := util.MapGetKeyAsType[int]("id", user.RawData)
	zauthName, err2 := util.MapGetKeyAsType[string]("fullName", user.RawData)
	zauthUsername, err3 := util.MapGetKeyAsType[string]("username", user.RawData)
	if err1 != nil || err2 != nil || err3 != nil {
		zap.S().Error(err1, err2, err3)
		return fiber.ErrInternalServerError
	}

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

	if dbUser.ID == 0 {
		// User is not in the databank yet
		// This means the user was never a board member (but somehow has admin)
		// Or the new board members haven't been scraped yet
		// Or the user's zauth fullname doesn't match with the one on the website
		// In any case, access denied
		return fiber.ErrUnauthorized
	}

	session, err := goth_fiber.SessionStore.Get(c)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	session.Set("memberID", dbUser.ID)
	if err = session.Save(); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return c.Redirect(r.redirectURL)
}

func (r *AuthRouter) logoutHandler(c *fiber.Ctx) error {
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
