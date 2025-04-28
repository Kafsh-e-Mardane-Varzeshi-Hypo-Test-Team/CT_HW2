// cmd/api/main.go
package main

import (
	"log"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/api"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/execute", api.ExecuteCode)

	if err := r.Run(":9090"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
