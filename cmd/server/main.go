package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/database/generated"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/config"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/database"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbPool, err := db.NewDBPool(ctx, cfg.Database.ConnectionString())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	queries := db.NewQuerier(dbPool)
	queries.UpdateUser(ctx, generated.UpdateUserParams{
		ID:       1,
		Username: "Arash Mohseni",
		EncryptedPassword: "I'm GOOD",
		Role: "admin"})
	fmt.Println(queries.ListUsers(ctx))
}