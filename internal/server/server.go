package server

import (
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/config"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/handlers"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/middlewares"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/services"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Engine   *gin.Engine
	Configs  *config.Config
	Database *services.DBService
}

func NewServer(engine *gin.Engine, configs *config.Config, database *services.DBService) *Server {
	return &Server{
		Engine:   engine,
		Configs:  configs,
		Database: database,
	}
}

func (s *Server) Start() {
	s.registerRoutes()
	s.Engine.Run(s.Configs.Server.Address())
}

func (s *Server) registerRoutes() {
	service := services.NewService(s.Configs, s.Database)
	middleware := middlewares.NewMiddleware(s.Configs, s.Database)
	handler := handlers.NewHandler(service)

	r := s.Engine

	r.Use(middleware.AuthMiddleware())

	r.GET("/", handler.IndexPage)

	r.GET("/login", handler.LoginPage)
	r.POST("/login", handler.LoginPost)
	r.GET("/signup", handler.SignupPage)
	r.POST("/signup", handler.SignupPost)
	r.POST("/logout", handler.Logout)

	r.GET("/profile/:username", handler.ProfilePage)
	r.POST("/demote-user", handler.DemoteUser)
	r.POST("/promote-user", handler.PromoteUser)

	r.GET("/problemset", handler.ProblemsetPage)

	r.GET("/addedproblems", handler.AddedProblemsPage)
	r.POST("/draft-problem", handler.DraftProblem)
	r.POST("/publish-problem", handler.PublishProblem)

	r.GET("/problem/:id", handler.ProblemPage)
	r.GET("/newproblem", handler.NewProblemPage)
	r.POST("/newproblem", handler.NewProblemPost)

	r.GET("/editproblem/:id", handler.EditProblemPage)
	r.POST("/editproblem", handler.EditProblemPost)

	r.GET("/submit", handler.SubmitPage)
	r.GET("/submit/:id", handler.SubmitPage)
	r.POST("/submit", handler.SubmitPost)
	r.GET("/submissions", handler.SubmissionsPage)
}
