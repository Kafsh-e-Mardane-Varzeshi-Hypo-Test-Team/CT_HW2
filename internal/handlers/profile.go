package handlers

import (
	"fmt"
	"strconv"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/services"
	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {
	// Extract user ID from the context
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// parse int
	fmt.Println(userID)
	userIDInt, ok := strconv.Atoi(fmt.Sprint(userID))
	if ok != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	// Fetch user profile from the database (mocked here)
	profile, err := services.GetUserProfile(userIDInt)
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not fetch profile"})
		return
	}

	c.JSON(200, profile)
}
