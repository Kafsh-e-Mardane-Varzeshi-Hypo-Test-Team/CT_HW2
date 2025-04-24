package middlewares

import (
	"net/http"
	"strconv"

	jwt "github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/pkg"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("session_token")
		if err != nil || token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		claims, valid := jwt.ValidateToken(token, m.Configs.JWT.SecretKey)
		if !valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		userId := claims["userId"]
		if userId == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token does not contain userId"})
			c.Abort()
			return
		}
		userIdInt, err := strconv.Atoi(userId.(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		user, err := m.Queries.GetUserById(c, int32(userIdInt))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("user", user)

		c.Next()
	}
}
