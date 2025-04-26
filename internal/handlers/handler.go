package handlers

import "github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/services"

type Handler struct {
	Service *services.Service
}

type User struct {
	ID       int32
	Username string
	IsAdmin  bool
}

type Problem struct {
	ID            int32
	Title         string
	Owner         string
	Status        string
	TimeLimitMs   int32
	MemoryLimitMb int32
	Statement     string
}

type Submission struct {
	ID      int32
	When    string
	Problem Problem
	Status  string
	Time    int
	Memory  int
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{
		Service: service,
	}
}
