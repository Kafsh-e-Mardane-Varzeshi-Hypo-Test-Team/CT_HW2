package services

import (
	"context"
	"fmt"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/database"
	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/database/generated"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBService struct {
	DB      *pgxpool.Pool
	Queries *generated.Queries
}

func InitDB(ctx context.Context, dbConStr string) (*DBService, error) {
	dbPool, err := database.NewDBPool(ctx, dbConStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	queries := database.NewQuerier(dbPool)
	return &DBService{
		DB:      dbPool,
		Queries: queries,
	}, nil
}
