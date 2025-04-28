package handlers

import (
	"io"
	"net/http"
	"strconv"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/database/generated"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	submissionsPageSize int = 10
)

func (h *Handler) SubmissionsPage(c *gin.Context) {
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
		c.Redirect(http.StatusFound, "/submissions?page=1")
		return
	}

	var submissionCnt int64

	submissionCnt, err = h.Service.Database.Queries.GetUserSubmissionsCount(c.Request.Context(), pgtype.Int4{
		Int32: user["ID"].(int32),
		Valid: true,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch submission count",
		})
		return
	}

	totalPages := max((submissionCnt+int64(submissionsPageSize)-1)/int64(submissionsPageSize), 1)

	// TODO: error page
	if (submissionCnt == 0 && currentPage != 1) || currentPage < 1 {
		c.Redirect(http.StatusFound, "/submissions?page=1")
		return
	} else if int64(currentPage) > totalPages {
		c.Redirect(http.StatusFound, "/submissions?page="+strconv.FormatInt(totalPages, 10))
		return
	}

	intervalStart := (currentPage - 1) * submissionsPageSize

	var submissions []generated.Submission

	submissions, err = h.Service.Database.Queries.ListUserSubmissions(c.Request.Context(), generated.ListUserSubmissionsParams{
		UserID: pgtype.Int4{Int32: user["ID"].(int32), Valid: true},
		Limit:  int32(submissionsPageSize),
		Offset: int32(intervalStart),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch submissions",
		})
		return
	}

	submissionsList := make([]Submission, 0, submissionsPageSize)

	for _, submission := range submissions {
		problem, err := h.Service.Database.Queries.GetProblemById(c.Request.Context(), submission.ProblemID.Int32)

		if err != nil {
			problem = generated.Problem{}
		}

		submissionsList = append(submissionsList, Submission{
			ID:   submission.ID,
			When: submission.SubmittedAt.Time.String(),
			Problem: Problem{
				ID:          problem.ID,
				Title:       problem.Title,
				Status:      toTitle(string(problem.Status)),
				TimeLimit:   problem.TimeLimitMs,
				MemoryLimit: problem.MemoryLimitMb,
				Statement:   problem.Statement,
				Input:       problem.SampleInput.String,
				Output:      problem.SampleOutput.String,
			},
			Status: toTitle(string(submission.Status)),
			Time:   int(submission.ExecutionTimeMs.Int32),
			Memory: int(submission.MemoryUsedMb.Int32),
		})
	}

	data["Submissions"] = submissionsList
	data["CurrentPage"] = currentPage
	data["TotalPages"] = totalPages

	c.HTML(http.StatusOK, "submissions.html", data)
	c.Abort()
}

func (h *Handler) SubmitPage(c *gin.Context) {
	user, exists := c.Get("User")

	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	problemID := c.Param("id")
	c.HTML(http.StatusOK, "submit.html", gin.H{
		"User": user,
		"ID":   problemID,
	})
}

func (h *Handler) SubmitPost(c *gin.Context) {
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

	problemID := c.PostForm("id")
	language := c.PostForm("language")
	method := c.PostForm("method")

	if language == "" || problemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Language and problem ID are required",
		})
		return
	}

	uploadedFile, err := c.FormFile("file")
	if method == "file" && err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File is required",
		})
		return
	}

	var code string
	if method == "file" {
		f, err := uploadedFile.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to open uploaded file",
			})
			return
		}
		defer f.Close()

		content, err := io.ReadAll(f)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to read uploaded file",
			})
			return
		}

		code = string(content)
	} else {
		code = c.PostForm("code")
		if method == "code" && code == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Code is required",
			})
			return
		}
	}

	problemIDInt, err := strconv.Atoi(problemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid problem ID",
		})
		return
	}

	problem, err := h.Service.Database.Queries.GetProblemById(c.Request.Context(), int32(problemIDInt))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Problem not found",
		})
		return
	}

	if problem.Status != generated.ProblemStatusPublished {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Problem is not published",
		})
		return
	}

	_, err = h.Service.Database.Queries.CreateSubmission(c.Request.Context(), generated.CreateSubmissionParams{
		UserID:     pgtype.Int4{Int32: user["ID"].(int32), Valid: true},
		ProblemID:  pgtype.Int4{Int32: problem.ID, Valid: true},
		SourceCode: code,
		Status:     generated.SubmissionStatusPending,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to submit code",
		})
		return
	}

	c.Redirect(http.StatusFound, "/submissions?page=1")
}
