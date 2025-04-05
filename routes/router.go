package routes

import (
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes initializes all API routes
func RegisterRoutes(r *gin.Engine) {
	// TOOD: CHANGE THIS
	r.Use(handlers.AuthMiddleware())

	r.GET("/", handlers.IndexPage)
	r.GET("/login", handlers.LoginPage)
	r.POST("/login", handlers.LoginHandler)
	r.GET("/signup", handlers.SignupPage)
	// r.POST("/signup", handlers.signupHandler)
	r.GET("/profile/:username", handlers.ProfilePage)
	// r.POST("/demote-user")
	// r.POST("/promote-user")
	r.GET("/problemset", handlers.ProblemsetPage)
	r.GET("/submit/:id", handlers.SubmitPage)
	r.GET("/submit", handlers.SubmitPage)
	// r.POST("/submit")
	r.GET("/submissions", handlers.SubmissionsPage)
	r.GET("/addedproblems", handlers.AddedProblemsPage)
	// r.POST("/draft-problem")
	// r.POST("/publish-problem")
	r.GET("/problem/:id", handlers.ProblemPage)
	r.GET("/newproblem", handlers.NewProblemPage)
	// r.POST("/newproblem")
	r.GET("/editproblem/:id", handlers.EditProblemPage)
	// r.POST("/editproblem")
}
