package api

type ExecuteRequest struct {
	Code           string `json:"code" binding:"required"`
	Input          string `json:"input" binding:"required"`
	ExpectedOutput string `json:"expected_output" binding:"required"`
	MemoryLimit    int    `json:"memory_limit" binding:"required"`
	TimeLimit      int    `json:"timeout" binding:"required"`
}

type ExecuteResponse struct {
	Status     string `json:"status"`
	TimeUsed   int    `json:"time_used"`
	MemoryUsed int    `json:"memory_used"`
	Error      string `json:"error,omitempty"`
}
