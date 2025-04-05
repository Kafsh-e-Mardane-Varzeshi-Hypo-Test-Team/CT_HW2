package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/config"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/database"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/database/generated"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/web/renderer"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbPool, err := database.NewDBPool(ctx, config.Database.ConnectionString())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	queries := database.NewQuerier(dbPool)
	queries.UpdateUser(ctx, generated.UpdateUserParams{
		ID:                1,
		Username:          "Arash Mohseni",
		EncryptedPassword: "I'm GOD",
		Role:              "admin"})
	fmt.Println(queries.ListUsers(ctx))

	// Init Server
	path := "internal/web/templates"

	r := gin.Default()
	r.Static("/static", "./static")
	r.HTMLRender = renderer.LoadTemplates(path)

	routes.RegisterRoutes(r)

	r.Run(config.Server.Address())
}
