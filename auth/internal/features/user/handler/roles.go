package handler

import (
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"

	"auth/internal/app"
	"auth/internal/shared/strutil"
)

func GetRoles(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		targetUserIDParam := c.Param("id")
		targetUserID, err := strutil.StringToUint(targetUserIDParam)
		if err != nil {
			zap.S().Errorf("Failed to parse target user ID:\n%f", err)

			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		roles, err := app.UserModule.Roles.GetRoles(targetUserID)
		if err != nil {
			zap.S().Errorf("Failed to get roles:\n%f", err)

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to get roles"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"roles": roles,
		})
	}
}

type GrantRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

func GrantRole(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		targetUserIDParam := c.Param("id")
		targetUserID, err := strutil.StringToUint(targetUserIDParam)
		if err != nil {
			zap.S().Errorf("Failed to parse target user ID:\n%f", err)

			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		exists, err := app.UserModule.User.UserExists(targetUserID)
		if err != nil {
			zap.S().Errorf("Failed to check if user exists:\n%f", err)

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		if !exists {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "target user not found"})
			return
		}

		req := GrantRoleRequest{}
		if err = c.ShouldBindJSON(&req); err != nil {
			zap.S().Errorf("Failed to parse grant role request:\n%f", err)

			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		wasGranted, err := app.UserModule.Roles.GrantRole(targetUserID, req.Role)
		if err != nil {
			zap.S().Errorf("Failed to grant role:\n%f", err)

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to grant role"})
			return
		}
		if !wasGranted {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "user already has this role"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "role has been granted",
		})
	}
}

func RevokeRole(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		targetUserIDParam := c.Param("id")
		targetUserID, err := strutil.StringToUint(targetUserIDParam)
		if err != nil {
			zap.S().Errorf("Failed to parse target user ID:\n%f", err)

			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		exists, err := app.UserModule.User.UserExists(targetUserID)
		if err != nil {
			zap.S().Errorf("Failed to check if user exists:\n%f", err)

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		if !exists {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "target user not found"})
			return
		}

		targetRole := c.Param("role")
		if targetRole == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
			return
		}

		err = app.UserModule.Roles.RevokeRole(targetUserID, targetRole)
		if err != nil {
			zap.S().Errorf("Failed to revoke role:\n%f", err)

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to revoke role"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "role has been revoked",
		})
	}
}
