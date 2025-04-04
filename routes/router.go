package routes

import (
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/handlers"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/middlewares"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes initializes all API routes
func RegisterRoutes(r *gin.Engine) {
	base := r.Group("/")
	{
		base.POST("/register", handlers.RegisterUser)
		base.POST("/login", handlers.LoginUser)
		base.GET("/users", handlers.GetAllUsers)

		protected := base.Group("/")
		protected.Use(middlewares.AuthMiddleware()) // Protect routes
		{
			protected.GET("/profile", handlers.GetProfile)
		}
	}
}
