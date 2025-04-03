package handlers

import (
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/services"
	"github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) {
	// Fetch all users from the database
	users := services.GetAllUsers()

	c.JSON(200, users)
}
