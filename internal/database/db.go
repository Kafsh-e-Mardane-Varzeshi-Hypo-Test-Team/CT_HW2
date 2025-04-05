package database

import (
	"context"
	"fmt"

	"github.com/Kafsh-e-Mardane-Varzeshi-Hypo-Test-Team/CT_HW2/internal/database/generated"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewDBPool creates a new connection pool to the PostgreSQL database
func NewDBPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	dbPool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Verify connection
	if err := dbPool.Ping(ctx); err != nil {
		dbPool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return dbPool, nil
}

// NewQuerier creates a new Querier with the provided connection pool
func NewQuerier(dbPool *pgxpool.Pool) *generated.Queries {
	return generated.New(dbPool)
}
