package handlers

import (
	"net/http"
	"strconv"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/database/generated"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
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
	user, exists := c.Get("User")

	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.HTML(http.StatusOK, "new_problem.html", gin.H{
		"User": user,
	})
}

func (h *Handler) NewProblemPost(c *gin.Context) {
	data := gin.H{}
	userCache, exists := c.Get("User")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	user, ok := userCache.(gin.H)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user",
		})
		return
	}

	data["User"] = user
	title := c.PostForm("title")
	statement := c.PostForm("statement")
	timeLimit := c.PostForm("time")
	memoryLimit := c.PostForm("memory")
	input := c.PostForm("input")
	output := c.PostForm("output")

	if title == "" || statement == "" || timeLimit == "" || memoryLimit == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "All fields are required",
		})
		return
	}

	timeLimitInt, err := strconv.Atoi(timeLimit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid time limit",
		})
		return
	}
	memoryLimitInt, err := strconv.Atoi(memoryLimit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid memory limit",
		})
		return
	}

	problem, err := h.Service.Database.Queries.CreateProblem(c, generated.CreateProblemParams{
		Title:         title,
		Statement:     statement,
		TimeLimitMs:   int32(timeLimitInt),
		MemoryLimitMb: int32(memoryLimitInt),
		SampleInput:   pgtype.Text{String: input, Valid: true},
		SampleOutput:  pgtype.Text{String: output, Valid: true},
		OwnerID:       user["ID"].(int32),
		Status:        generated.ProblemStatusDraft,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create problem",
		})
		return
	}

	c.Redirect(http.StatusFound, "/problem/"+strconv.Itoa(int(problem.ID)))
}

func (h *Handler) EditProblemPage(c *gin.Context) {
	data := gin.H{}
	userCache, exists := c.Get("User")

	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	user, ok := userCache.(gin.H)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user",
		})
		return
	}

	data["User"] = user

	// fetch user from database
	userMakingChange, err := h.Service.Database.Queries.GetUserById(c, user["ID"].(int32))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user",
		})
		return
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

	if problem.OwnerID != userMakingChange.ID && userMakingChange.Role != generated.UserRoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not allowed to edit this problem",
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
		Input:       problem.SampleInput.String,
		Output:      problem.SampleOutput.String,
	}

	c.HTML(http.StatusOK, "edit_problem.html", data)
}

func (h *Handler) EditProblemPost(c *gin.Context) {
	data := gin.H{}
	userCache, exists := c.Get("User")

	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	user, ok := userCache.(gin.H)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user",
		})
		return
	}

	data["User"] = user

	// fetch user from database
	userMakingChange, err := h.Service.Database.Queries.GetUserById(c, user["ID"].(int32))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user",
		})
		return
	}

	id_str := c.PostForm("id")
	id, err := strconv.Atoi(id_str)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid problem ID",
		})
		return
	}

	data["User"] = user
	title := c.PostForm("title")
	statement := c.PostForm("statement")
	timeLimit := c.PostForm("time")
	memoryLimit := c.PostForm("memory")
	input := c.PostForm("input")
	output := c.PostForm("output")

	if title == "" || statement == "" || timeLimit == "" || memoryLimit == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "All fields are required",
		})
		return
	}

	timeLimitInt, err := strconv.Atoi(timeLimit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid time limit",
		})
		return
	}
	memoryLimitInt, err := strconv.Atoi(memoryLimit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid memory limit",
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

	if problem.OwnerID != userMakingChange.ID && userMakingChange.Role != generated.UserRoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not allowed to edit this problem",
		})
		return
	}

	updatedProblem, err := h.Service.Database.Queries.UpdateProblem(c, generated.UpdateProblemParams{
		Title:         title,
		Statement:     statement,
		TimeLimitMs:   int32(timeLimitInt),
		MemoryLimitMb: int32(memoryLimitInt),
		SampleInput:   pgtype.Text{String: input, Valid: true},
		SampleOutput:  pgtype.Text{String: output, Valid: true},
		Status:        generated.ProblemStatusDraft,
		ID:            int32(id),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create problem",
		})
		return
	}

	c.Redirect(http.StatusFound, "/problem/"+strconv.Itoa(int(updatedProblem.ID)))
}
