package middlewares

import (
	"net/http"

	jwt "github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/pkg"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		claims, valid := jwt.ValidateToken(token)
		if !valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		userId, exists := claims["userId"].(string)
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "userId not found"})
			c.Abort()
			return
		}

		// TODO: query the database to get the user role

		// TODO: check if token is expired (both cases of time and change of role or removal of user)

		// Store user data in the request context
		c.Set("userId", userId)
		// c.Set("role", role)

		c.Next()
	}
}
