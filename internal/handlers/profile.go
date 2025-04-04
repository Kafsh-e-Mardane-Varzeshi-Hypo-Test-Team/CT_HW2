package handlers

import (
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/services"
	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {
	user, ok := getUser(c)
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// Fetch user profile from the database (mocked here)
	profile, err := services.GetUserProfile(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not fetch profile"})
		return
	}

	c.JSON(200, profile)
}
