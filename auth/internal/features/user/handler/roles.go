package handler

import (
	"log"
	"net/http"

	"github.com/DangeL187/erax"
	"github.com/gin-gonic/gin"

	"auth/internal/app"
	"auth/internal/shared/strutil"
)

func GetRoles(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		targetUserIDParam := c.Param("id")
		targetUserID, err := strutil.StringToUint(targetUserIDParam)
		if err != nil {
			err = erax.Wrap(err, "failed to parse target user ID")
			log.Printf("\n%f\n", err)

			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		roles, err := app.UserModule.Roles.GetRoles(targetUserID)
		if err != nil {
			err = erax.Wrap(err, "failed to get roles")
			log.Printf("\n%f\n", err)

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
			err = erax.Wrap(err, "failed to parse target user ID")
			log.Printf("\n%f\n", err)

			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		exists, err := app.UserModule.User.UserExists(targetUserID)
		if err != nil {
			err = erax.Wrap(err, "failed to check if user exists")
			log.Printf("\n%f\n", err)

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		if !exists {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "target user not found"})
			return
		}

		req := GrantRoleRequest{}
		if err = c.ShouldBindJSON(&req); err != nil {
			err = erax.Wrap(err, "failed to parse grant role request")
			log.Printf("\n%f\n", err)

			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		wasGranted, err := app.UserModule.Roles.GrantRole(targetUserID, req.Role)
		if err != nil {
			err = erax.Wrap(err, "failed to grant role")
			log.Printf("\n%f\n", err)

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
			err = erax.Wrap(err, "failed to parse target user ID")
			log.Printf("\n%f\n", err)

			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		exists, err := app.UserModule.User.UserExists(targetUserID)
		if err != nil {
			err = erax.Wrap(err, "failed to check if user exists")
			log.Printf("\n%f\n", err)

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
			err = erax.Wrap(err, "failed to revoke role")
			log.Printf("\n%f\n", err)

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to revoke role"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "role has been revoked",
		})
	}
}
