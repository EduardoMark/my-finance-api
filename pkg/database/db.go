package database

import (
	"database/sql"
	"fmt"

	"github.com/EduardoMark/my-finance-api/pkg/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectDatabase(cfg *config.Env) (*sql.DB, error) {
	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBTimezone,
	)

	DB, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	err = DB.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return DB, nil
}
