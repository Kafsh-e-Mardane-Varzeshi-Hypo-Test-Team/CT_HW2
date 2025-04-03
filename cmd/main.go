package main

import (
	"log"
	"os"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/config"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// TODO: docker-compose up
	config.LoadEnv()

	r := gin.Default()

	routes.RegisterRoutes(r)

	if port, exists := os.LookupEnv("PROJECT_PORT"); exists {
		log.Println("Starting server on port", port)
		r.Run(":" + port)
	}
}
