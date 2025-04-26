package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) IndexPage(c *gin.Context) {
	user, exists := c.Get("User")

	if !exists {
		c.HTML(http.StatusOK, "index.html", nil)
	} else {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"User": user,
		})
		return
	}
}
