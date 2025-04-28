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
	ID          int32
	Title       string
	Owner       string
	Status      string
	TimeLimit   int32
	MemoryLimit int32
	Statement   string
	Input       string
	Output      string
}

type Submission struct {
	ID      int32
	When    string
	Problem Problem
	Status  string
	Time    int
	Memory  int
}

type Profile struct {
	ID                    int32
	Username              string
	IsAdmin               bool
	TotalSubmissions      int
	SuccessfulSubmissions int
	Submissions           []Submission
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{
		Service: service,
	}
}
