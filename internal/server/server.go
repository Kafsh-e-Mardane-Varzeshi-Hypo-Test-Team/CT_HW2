package server

import (
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/config"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/database/generated"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/handlers"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/middlewares"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/services"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Engine  *gin.Engine
	Configs *config.Config
	Queries *generated.Queries
}

func NewServer(engine *gin.Engine, configs *config.Config, queries *generated.Queries) *Server {
	return &Server{
		Engine:  engine,
		Configs: configs,
		Queries: queries,
	}
}

func (s *Server) Start() {
	s.registerRoutes()
	s.Engine.Run(s.Configs.Server.Address())
}

func (s *Server) registerRoutes() {
	service := services.NewService(s.Configs, s.Queries)
	middleware := middlewares.NewMiddleware(s.Configs, s.Queries)
	handler := handlers.NewHandler(service)

	r := s.Engine

	r.Use(middleware.AuthMiddleware())

	r.GET("/", handler.IndexPage)

	r.GET("/login", handler.LoginGet)
	r.POST("/login", handler.LoginPost)
	r.GET("/signup", handler.SignupGet)
	r.POST("/signup", handler.SignupPost)
	r.POST("/logout", handler.Logout)

	r.GET("/profile/:username", handler.ProfilePage)
	// r.POST("/demote-user")
	// r.POST("/promote-user")
	r.GET("/problemset", handler.ProblemsetPage)
	r.GET("/submit/:id", handler.SubmitPage)
	r.GET("/submit", handler.SubmitPage)
	// r.POST("/submit")
	r.GET("/submissions", handler.SubmissionsPage)
	r.GET("/addedproblems", handler.AddedProblemsPage)
	// r.POST("/draft-problem")
	// r.POST("/publish-problem")
	r.GET("/problem/:id", handler.ProblemPage)
	r.GET("/newproblem", handler.NewProblemPage)
	// r.POST("/newproblem")
	r.GET("/editproblem/:id", handler.EditProblemPage)
	// r.POST("/editproblem")
}
