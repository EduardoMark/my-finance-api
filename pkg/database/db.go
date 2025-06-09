package database

import (
	"context"
	"fmt"

	"github.com/EduardoMark/my-finance-api/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDatabase(cfg *config.Env) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBTimezone,
	)

	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := dbpool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return dbpool, nil
}
