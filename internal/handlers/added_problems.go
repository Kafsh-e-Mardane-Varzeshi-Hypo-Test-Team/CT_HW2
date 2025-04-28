package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/database/generated"
	"github.com/gin-gonic/gin"
)

const (
	addedProblemsPageSize int = 10
)

func (h *Handler) AddedProblemsPage(c *gin.Context) {
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

	currentPage, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.Redirect(http.StatusFound, "/addedproblems?page=1")
		return
	}

	var problemCnt int64

	if user["IsAdmin"].(bool) {
		problemCnt, err = h.Service.Database.Queries.GetProblemsCount(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch problem count",
			})
			return
		}
	} else {
		problemCnt, err = h.Service.Database.Queries.GetUserProblemsCount(c.Request.Context(), user["ID"].(int32))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch problem count",
			})
			return
		}
	}

	totalPages := max((problemCnt+int64(addedProblemsPageSize)-1)/int64(addedProblemsPageSize), 1)

	// TODO: error page
	if (problemCnt == 0 && currentPage != 1) || currentPage < 1 {
		c.Redirect(http.StatusFound, "/addedproblems?page=1")
		return
	} else if int64(currentPage) > totalPages {
		fmt.Println("currentPage", currentPage, "totalPages", totalPages, "problemCnt", problemCnt)
		c.Redirect(http.StatusFound, "/addedproblems?page="+strconv.FormatInt(totalPages, 10))
		return
	}

	intervalStart := (currentPage - 1) * addedProblemsPageSize

	var problems []generated.Problem

	if user["IsAdmin"].(bool) {
		problems, err = h.Service.Database.Queries.ListProblems(c.Request.Context(), generated.ListProblemsParams{
			Limit:  int32(addedProblemsPageSize),
			Offset: int32(intervalStart),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch problems",
			})
			return
		}
	} else {
		problems, err = h.Service.Database.Queries.ListUserProblems(c.Request.Context(), generated.ListUserProblemsParams{
			OwnerID: user["ID"].(int32),
			Limit:   int32(addedProblemsPageSize),
			Offset:  int32(intervalStart),
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch problems",
			})
			return
		}
	}

	problemList := make([]Problem, 0, addedProblemsPageSize)

	// search for owner user, use map probably to do fewer queries
	// define map:

	var userMap map[int32]generated.User = make(map[int32]generated.User)

	for _, problem := range problems {
		ownerName := "Unknown"
		if owner, exists := userMap[problem.OwnerID]; exists {
			ownerName = owner.Username
		} else {
			owner, err := h.Service.Database.Queries.GetUserById(c.Request.Context(), problem.OwnerID)
			if err == nil {
				ownerName = owner.Username
				userMap[problem.OwnerID] = owner
			} else {
				fmt.Println("Error fetching owner:", err)
			}
		}

		problemList = append(problemList, Problem{
			ID:          problem.ID,
			Title:       problem.Title,
			Owner:       ownerName,
			Status:      toTitle(string(problem.Status)),
			TimeLimit:   problem.TimeLimitMs,
			MemoryLimit: problem.MemoryLimitMb,
			Statement:   problem.Statement,
		})
	}

	data["Problems"] = problemList
	data["CurrentPage"] = currentPage
	data["TotalPages"] = totalPages

	c.HTML(http.StatusOK, "added_problems.html", data)
	c.Abort()
}

// TODO: same function as published, refactor these
func (h *Handler) DraftProblem(c *gin.Context) {
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

	tx, err := h.Service.Database.DB.Begin(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to begin transaction",
		})
		return
	}
	defer tx.Rollback(c.Request.Context())

	userMakingChange, err := h.Service.Database.Queries.WithTx(tx).GetUserById(c.Request.Context(), user["ID"].(int32))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user",
		})
		return
	}

	if userMakingChange.Role != generated.UserRoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You do not have permission to perform this action",
		})
		return
	}

	problem_id_string := c.PostForm("problem_id")
	problem_id, err := strconv.Atoi(problem_id_string)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid problem ID",
		})
		return
	}

	problem, err := h.Service.Database.Queries.WithTx(tx).GetProblemById(c.Request.Context(), int32(problem_id))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch problem",
		})
		return
	}

	_, err = h.Service.Database.Queries.WithTx(tx).UpdateProblem(c.Request.Context(), generated.UpdateProblemParams{
		Title:         problem.Title,
		Statement:     problem.Statement,
		TimeLimitMs:   problem.TimeLimitMs,
		MemoryLimitMb: problem.MemoryLimitMb,
		SampleInput:   problem.SampleInput,
		SampleOutput:  problem.SampleOutput,
		Status:        generated.ProblemStatusDraft,
		ID:            problem.ID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to promote user",
		})
		return
	}

	err = tx.Commit(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to commit transaction",
		})
		return
	}

	c.Redirect(http.StatusFound, c.Request.Referer())
	c.Abort()
}

func (h *Handler) PublishProblem(c *gin.Context) {
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

	tx, err := h.Service.Database.DB.Begin(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to begin transaction",
		})
		return
	}
	defer tx.Rollback(c.Request.Context())

	userMakingChange, err := h.Service.Database.Queries.WithTx(tx).GetUserById(c.Request.Context(), user["ID"].(int32))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user",
		})
		return
	}

	if userMakingChange.Role != generated.UserRoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You do not have permission to perform this action",
		})
		return
	}

	problem_id_string := c.PostForm("problem_id")
	problem_id, err := strconv.Atoi(problem_id_string)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid problem ID",
		})
		return
	}

	problem, err := h.Service.Database.Queries.WithTx(tx).GetProblemById(c.Request.Context(), int32(problem_id))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch problem",
		})
		return
	}

	_, err = h.Service.Database.Queries.WithTx(tx).UpdateProblem(c.Request.Context(), generated.UpdateProblemParams{
		Title:         problem.Title,
		Statement:     problem.Statement,
		TimeLimitMs:   problem.TimeLimitMs,
		MemoryLimitMb: problem.MemoryLimitMb,
		SampleInput:   problem.SampleInput,
		SampleOutput:  problem.SampleOutput,
		Status:        generated.ProblemStatusPublished,
		ID:            problem.ID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to promote user",
		})
		return
	}

	err = tx.Commit(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to commit transaction",
		})
		return
	}

	c.Redirect(http.StatusFound, c.Request.Referer())
	c.Abort()
}

func toTitle(s string) string {
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}
