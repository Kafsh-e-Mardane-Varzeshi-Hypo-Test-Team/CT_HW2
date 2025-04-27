package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/database/generated"
	"github.com/gin-gonic/gin"
)

const (
	problemsetPageSize int = 10
)

func (h *Handler) ProblemsetPage(c *gin.Context) {
	data := gin.H{}

	if user, exists := c.Get("User"); exists {
		data["User"] = user
	}

	currentPage, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.Redirect(http.StatusFound, "/problemset?page=1")
		return
	}

	problemCnt, err := h.Service.Database.Queries.GetPublishedProblemsCount(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch problem count",
		})
		return
	}

	totalPages := max((problemCnt+int64(problemsetPageSize)-1)/int64(problemsetPageSize), 1)

	// TODO: error page
	if (problemCnt == 0 && currentPage != 1) || currentPage < 1 {
		c.Redirect(http.StatusFound, "/problemset?page=1")
		return
	} else if int64(currentPage) > totalPages {
		fmt.Println("currentPage", currentPage, "totalPages", totalPages, "problemCnt", problemCnt)
		c.Redirect(http.StatusFound, "/problemset?page="+strconv.FormatInt(totalPages, 10))
		return
	}

	intervalStart := (currentPage - 1) * problemsetPageSize

	problems, err := h.Service.Database.Queries.ListPublishedProblems(c.Request.Context(), generated.ListPublishedProblemsParams{
		Limit:  int32(problemsetPageSize),
		Offset: int32(intervalStart),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch problems",
		})
		return
	}

	problemList := make([]Problem, 0, problemsetPageSize)

	for _, problem := range problems {
		problemList = append(problemList, Problem{
			ID:            problem.ID,
			Title:         problem.Title,
			Owner:         string(problem.OwnerID),
			Status:        toTitle(string(problem.Status)),
			TimeLimitMs:   problem.TimeLimitMs,
			MemoryLimitMb: problem.MemoryLimitMb,
			Statement:     problem.Statement,
		})
	}

	data["Problems"] = problemList
	data["CurrentPage"] = currentPage
	data["TotalPages"] = totalPages

	c.HTML(http.StatusOK, "problemset.html", data)
}
