package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ProblemPage(c *gin.Context) {
	data := gin.H{}

	if user, exists := c.Get("User"); exists {
		data["User"] = user
	}

	id_str := c.Param("id")
	id, err := strconv.Atoi(id_str)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid problem ID",
		})
		return
	}

	problem, err := h.Service.Database.Queries.GetProblemById(c, int32(id))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Problem not found",
		})
		return
	}

	ownerName := "Not Found"
	owner, err := h.Service.Database.Queries.GetUserById(c, problem.OwnerID)

	if err == nil {
		ownerName = owner.Username
	}

	data["Problem"] = Problem{
		ID:          problem.ID,
		Title:       problem.Title,
		Owner:       ownerName,
		Status:      toTitle(string(problem.Status)),
		TimeLimit:   problem.TimeLimitMs,
		MemoryLimit: problem.MemoryLimitMb,
		Statement:   problem.Statement,
	}

	c.HTML(http.StatusOK, "problem.html", data)
}

func (h *Handler) NewProblemPage(c *gin.Context) {
	c.HTML(http.StatusOK, "new_problem.html", gin.H{
		"User": admin,
	})
}

func (h *Handler) NewProblemPost(c *gin.Context) {
}

func (h *Handler) EditProblemPage(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	c.HTML(http.StatusOK, "edit_problem.html", gin.H{
		"Problem": problems[id-1],
		"User":    admin,
	})
}

func (h *Handler) EditProblemPost(c *gin.Context) {
}
