package handlers

import "github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/services"

type Handler struct {
	Service *services.Service
}

type User struct {
	ID       int
	Username string
	IsAdmin  bool
}

type Problem struct {
	ID          int
	Title       string
	Owner       string
	Status      string
	TimeLimit   int
	MemoryLimit int
	Statement   string
}

type Submission struct {
	ID      int
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
