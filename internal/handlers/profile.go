package handlers

import (
	"net/http"
	"strconv"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/database/generated"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	recentSubmissionsLimit = 5
)

func (h *Handler) ProfileGet(c *gin.Context) {
	user, exists := c.Get("User")

	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	profileUsername := c.Param("username")

	profileUser, err := h.Service.Database.Queries.GetUserByUsername(c.Request.Context(), profileUsername)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user",
		})
		return
	}
	profileUserStats, err := h.Service.Database.Queries.GetUserStatsById(c.Request.Context(), profileUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user stats",
		})
		return
	}

	recentSubmissions, err := h.Service.Database.Queries.ListUserSubmissions(c.Request.Context(), generated.ListUserSubmissionsParams{
		UserID: pgtype.Int4{Int32: profileUser.ID, Valid: true},
		Limit:  recentSubmissionsLimit,
		Offset: 0,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch recent submissions",
		})
		return
	}

	// ListUserSubmissions returns a joined result set of submissions and problems

	submissions := make([]Submission, 0, len(recentSubmissions))
	for _, submission := range recentSubmissions {
		problem := Problem{
			ID:            submission.ProblemID.Int32,
			Title:         submission.Title,
			Owner:         strconv.Itoa(int(submission.OwnerID)),
			Status:        string(submission.Status_2),
			TimeLimitMs:   submission.TimeLimitMs,
			MemoryLimitMb: submission.MemoryLimitMb,
			Statement:     submission.Statement,
		}

		submissions = append(submissions, Submission{
			ID:      submission.ID,
			When:    submission.SubmittedAt.Time.String(),
			Problem: problem,
			Status:  string(submission.Status),
			Time:    int(submission.ExecutionTimeMs.Int32),
			Memory:  int(submission.MemoryUsedMb.Int32),
		})
	}

	c.HTML(http.StatusOK, "profile.html", gin.H{
		"User": user,
		"Profile": Profile{
			ID:                    profileUser.ID,
			Username:              profileUser.Username,
			IsAdmin:               profileUser.Role == generated.UserRoleAdmin,
			TotalSubmissions:      int(profileUserStats.TotalSubmissions),
			SuccessfulSubmissions: int(profileUserStats.TotalAccepted),
			Submissions:           submissions,
		},
	})
}

func (h *Handler) DemoteUser(c *gin.Context) {
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

	user_id_string := c.PostForm("user_id")
	user_id, err := strconv.Atoi(user_id_string)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	if user_id == int(userMakingChange.ID) {
		c.Redirect(http.StatusFound, "/profile/"+userMakingChange.Username)
		return
	}

	userToUpdate, err := h.Service.Database.Queries.WithTx(tx).GetUserById(c.Request.Context(), int32(user_id))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user",
		})
		return
	}

	updatedUser, err := h.Service.Database.Queries.WithTx(tx).UpdateUser(c.Request.Context(), generated.UpdateUserParams{
		Username:          userToUpdate.Username,
		EncryptedPassword: userToUpdate.EncryptedPassword,
		Role:              generated.UserRoleNormal,
		ID:                int32(user_id),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to demote user",
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

	c.Redirect(http.StatusFound, "/profile/"+updatedUser.Username)
	c.Abort()
}

func (h *Handler) PromoteUser(c *gin.Context) {
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

	user_id_string := c.PostForm("user_id")
	user_id, err := strconv.Atoi(user_id_string)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	if user_id == int(userMakingChange.ID) {
		c.Redirect(http.StatusFound, "/profile/"+userMakingChange.Username)
		return
	}

	userToUpdate, err := h.Service.Database.Queries.WithTx(tx).GetUserById(c.Request.Context(), int32(user_id))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user",
		})
		return
	}

	updatedUser, err := h.Service.Database.Queries.WithTx(tx).UpdateUser(c.Request.Context(), generated.UpdateUserParams{
		Username:          userToUpdate.Username,
		EncryptedPassword: userToUpdate.EncryptedPassword,
		Role:              generated.UserRoleAdmin,
		ID:                int32(user_id),
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

	c.Redirect(http.StatusFound, "/profile/"+updatedUser.Username)
	c.Abort()
}
