package routes

import (
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/handlers"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/middlewares"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes initializes all API routes
func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/register", handlers.RegisterUser)
		api.POST("/login", handlers.LoginUser)

		protected := api.Group("/")
		protected.Use(middlewares.AuthMiddleware()) // Protect routes
		// {
		// 	protected.GET("/profile", handlers.GetProfile)
		// }
	}
}
