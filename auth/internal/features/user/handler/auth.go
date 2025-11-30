package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"auth/internal/app"
	"auth/internal/features/user/usecase"
	"auth/internal/infra/http/handlerutil"
	"auth/internal/shared/auth"
)

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginUserRequest
		if !handlerutil.BindJSON(c, &req, "failed to parse login user request") {
			return
		}

		input := usecase.LoginInput{
			Email:    req.Email,
			Password: req.Password,
		}

		accessToken, err := app.UserModule.Auth.Login(app.Config.UserAccessTokenTTL, input)
		if err != nil {
			handlerutil.HandleError(c, err, "failed to login user", map[error]handlerutil.ErrorResponse{
				auth.ErrInvalidCredentials: {http.StatusUnauthorized, "Invalid credentials"},
			})
			return
		}

		handlerutil.RespondAuth(c, handlerutil.TokenPair{
			AccessToken: accessToken,
		}, handlerutil.TokenConfig{
			AccessTTL: app.Config.UserAccessTokenTTL,
		}, "Logged in successfully")
	}
}

type RegisterUserRequest struct {
	Email    string `json:"email" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterUserRequest
		if !handlerutil.BindJSON(c, &req, "failed to parse register user request") {
			return
		}

		input := usecase.RegisterInput{
			Email:    req.Email,
			FullName: req.FullName,
			Password: req.Password,
		}

		accessToken, err := app.UserModule.Auth.Register(app.Config.UserAccessTokenTTL, input)
		if err != nil {
			handlerutil.HandleError(c, err, "failed to register user", map[error]handlerutil.ErrorResponse{
				auth.ErrInvalidCredentials: {http.StatusUnauthorized, "Invalid credentials"},
			})
			return
		}

		handlerutil.RespondAuth(c, handlerutil.TokenPair{
			AccessToken: accessToken,
		}, handlerutil.TokenConfig{
			AccessTTL: app.Config.UserAccessTokenTTL,
		}, "User registered successfully")
	}
}
