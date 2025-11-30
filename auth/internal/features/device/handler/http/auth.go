package http

import (
	"errors"
	"log"
	"net/http"

	"github.com/DangeL187/erax"
	"github.com/gin-gonic/gin"

	"auth/internal/app"
	"auth/internal/features/device/usecase"
	"auth/internal/infra/http/handlerutil"
	"auth/internal/shared/auth"
)

type AuthDeviceRequest struct {
	DeviceID string `json:"device_id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AuthDeviceRequest
		if !handlerutil.BindJSON(c, &req, "failed to parse login device request") {
			return
		}

		input := usecase.AuthInput{
			DeviceID: req.DeviceID,
			Password: req.Password,
		}

		accessToken, refreshToken, err := app.DeviceModule.Auth.Login(
			app.Config.DeviceAccessTokenTTL,
			app.Config.DeviceRefreshTokenTTL,
			input,
		)
		if err != nil {
			handlerutil.HandleError(c, err, "failed to login device", map[error]handlerutil.ErrorResponse{
				auth.ErrInvalidCredentials: {http.StatusUnauthorized, "Invalid credentials"},
			})
			return
		}

		handlerutil.RespondAuth(c, handlerutil.TokenPair{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}, handlerutil.TokenConfig{
			AccessTTL:   app.Config.DeviceAccessTokenTTL,
			RefreshTTL:  app.Config.DeviceRefreshTokenTTL,
			WithRefresh: true,
		}, "Logged in successfully")
	}
}

func Register(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AuthDeviceRequest
		if !handlerutil.BindJSON(c, &req, "failed to parse register device request") {
			return
		}

		input := usecase.AuthInput{
			DeviceID: req.DeviceID,
			Password: req.Password,
		}

		accessToken, refreshToken, err := app.DeviceModule.Auth.Register(
			app.Config.DeviceAccessTokenTTL,
			app.Config.DeviceRefreshTokenTTL,
			input,
		)
		if err != nil {
			err = erax.Wrap(err, "failed to register device")
			log.Printf("\n%f\n", err)

			if errors.Is(err, auth.ErrInvalidCredentials) {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid credentials"})
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		handlerutil.RespondAuth(c, handlerutil.TokenPair{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}, handlerutil.TokenConfig{
			AccessTTL:   app.Config.DeviceAccessTokenTTL,
			RefreshTTL:  app.Config.DeviceRefreshTokenTTL,
			WithRefresh: true,
		}, "Device registered successfully")
	}
}

func Refresh(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			err = erax.Wrap(err, "failed to find refresh_token cookie")
			log.Printf("\n%f\n", err)

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		accessToken, err := app.DeviceModule.Auth.Refresh(app.Config.DeviceAccessTokenTTL, refreshToken)
		if err != nil {
			err = erax.Wrap(err, "failed to refresh access token")
			log.Printf("\n%f\n", err)

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}

		handlerutil.RespondAuth(c, handlerutil.TokenPair{
			AccessToken: accessToken,
		}, handlerutil.TokenConfig{
			AccessTTL: app.Config.DeviceAccessTokenTTL,
		}, "Access token refreshed successfully")
	}
}
