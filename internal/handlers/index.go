package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) IndexPage(c *gin.Context) {
	data := gin.H{}

	if user, exists := c.Get("User"); exists {
		data["User"] = user
	}

	c.HTML(http.StatusOK, "index.html", data)
}
