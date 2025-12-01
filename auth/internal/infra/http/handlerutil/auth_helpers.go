package handlerutil

import (
	"errors"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Status  int
	Message string
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type TokenConfig struct {
	AccessTTL   time.Duration
	RefreshTTL  time.Duration
	WithRefresh bool
}

func BindJSON[T any](c *gin.Context, req *T, errorMessage string) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		zap.S().Errorf("%s:\n%f", errorMessage, err)

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return false
	}

	return true
}

func RespondAuth(c *gin.Context, tokens TokenPair, cfg TokenConfig, message string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(cfg.AccessTTL.Seconds()),
	})

	if cfg.WithRefresh {
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "refresh_token",
			Value:    tokens.RefreshToken,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   int(cfg.RefreshTTL.Seconds()),
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": message})
}

func HandleError(c *gin.Context, err error, message string, mapping map[error]ErrorResponse) {
	zap.S().Errorf("%s:\n%f", message, err)

	for target, errorResponse := range mapping {
		if errors.Is(err, target) {
			c.AbortWithStatusJSON(errorResponse.Status, gin.H{"error": errorResponse.Message})
			return
		}
	}

	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}
