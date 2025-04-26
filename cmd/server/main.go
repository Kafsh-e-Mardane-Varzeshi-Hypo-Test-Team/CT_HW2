package main

import (
	"context"
	"log"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/config"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/server"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/services"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/web/renderer"
	"github.com/gin-gonic/gin"
)

func main() {
	configs, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbService, err := services.InitDB(ctx, configs.Database.ConnectionString())
	if err != nil {
		log.Fatalf("DB init error: %v", err)
	}
	defer dbService.DB.Close()

	// Init Server
	path := "internal/web/templates"

	r := gin.Default()
	r.Static("/static", "./static")
	r.HTMLRender = renderer.LoadTemplates(path)

	server := server.NewServer(r, configs, dbService)
	server.Start()
}
