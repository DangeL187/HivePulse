package routes

import (
	"github.com/gin-gonic/gin"

	"auth/internal/app"
	deviceHandler "auth/internal/features/device/handler/http"
	userHandler "auth/internal/features/user/handler"
	"auth/internal/features/user/middleware"
)

func SetupRoutes(router *gin.Engine, app *app.App) {
	router.POST(
		"/users/login",
		userHandler.Login(app),
	)

	router.POST(
		"/users/register",
		userHandler.Register(app),
	)

	router.GET(
		"/users/:id/roles",
		middleware.Auth(app),
		middleware.EnsureSelfByID(app),
		userHandler.GetRoles(app),
	)

	router.POST(
		"/users/:id/roles",
		middleware.Auth(app),
		middleware.UserHasPermission(app, "user", "grant_role"),
		userHandler.GrantRole(app),
	)

	router.DELETE(
		"/users/:id/roles/:role",
		middleware.Auth(app),
		middleware.UserHasPermission(app, "user", "revoke_role"),
		userHandler.RevokeRole(app),
	)

	router.GET(
		"/users/lookup",
		middleware.Auth(app),
		middleware.UserHasPermission(app, "user", "view"),
		userHandler.Lookup(app),
	)

	router.POST(
		"/devices/login",
		deviceHandler.Login(app),
	)

	router.POST(
		"/devices/register",
		middleware.Auth(app),
		middleware.UserHasPermission(app, "device", "register"),
		deviceHandler.Register(app),
	)

	router.POST(
		"/devices/refresh",
		deviceHandler.Refresh(app),
	)
}
