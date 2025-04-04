package handlers

import (
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/models"
	"github.com/gin-gonic/gin"
)

func getUser(c *gin.Context) (*models.User, bool) {
	userVar, exists := c.Get("user")
	if !exists {
		return nil, false
	}
	user, ok := userVar.(*models.User)
	if !ok {
		return nil, false
	}

	return user, true
}
