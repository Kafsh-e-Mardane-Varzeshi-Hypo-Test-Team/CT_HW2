package middlewares

import (
	"net/http"
	"strconv"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/database/generated"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/handlers"
	jwt "github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/pkg"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("session_token")
		if err != nil || token == "" {
			handlers.ClearSessionCookie(c)
			c.Next()
			return
		}

		claims, valid := jwt.ValidateToken(token, m.Configs.JWT.SecretKey)
		if !valid {
			handlers.ClearSessionCookie(c)
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		userId := claims["userId"]
		if userId == nil {
			handlers.ClearSessionCookie(c)
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		userIdInt, err := strconv.Atoi(userId.(string))
		if err != nil {
			handlers.ClearSessionCookie(c)
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// TODO: session + refresh token, if refresh token is expired then query
		user, err := m.Queries.GetUserById(c, int32(userIdInt))
		if err != nil {
			handlers.ClearSessionCookie(c)
			c.Redirect(http.StatusUnauthorized, "/login")
			c.Abort()
			return
		}

		c.Set("User", gin.H{
			"ID":       user.ID,
			"Username": user.Username,
			"IsAdmin":  user.Role == generated.UserRoleAdmin,
		})

		c.Next()
	}
}
