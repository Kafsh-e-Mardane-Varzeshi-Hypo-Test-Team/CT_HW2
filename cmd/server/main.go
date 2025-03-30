package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/database/generated"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/config"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/database"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println(cfg.Database.ConnectionString())
	// Connect to database
	dbPool, err := db.NewDBPool(ctx, cfg.Database.ConnectionString())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Create database queries
	queries := db.NewQuerier(dbPool)

	queries.UpdateUser(ctx, generated.UpdateUserParams{
		ID:       1,
		Username: "Arash Mohseni",
		EncryptedPassword: "I'm GOOD",
		Role: "admin"})

	fmt.Println(queries.ListUsers(ctx))

	// Create shutdown context with timeout
	_, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

}