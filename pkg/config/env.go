package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBTimezone string
	JWTSecret  string
}

func LoadEnv() (*Env, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error on loading enviroments: %w", err)
	}

	return &Env{
		Port:       gentEnv("PORT", ""),
		DBHost:     gentEnv("DB_HOST", ""),
		DBPort:     gentEnv("DB_PORT", ""),
		DBUser:     gentEnv("DB_USER", ""),
		DBPassword: gentEnv("DB_PASSWORD", ""),
		DBName:     gentEnv("DB_NAME", ""),
		DBTimezone: gentEnv("DB_TIMEZONE", ""),
		JWTSecret:  gentEnv("JWT_SECRET", "secret"),
	}, nil
}

func gentEnv(key, fallback string) string {
	value, exist := os.LookupEnv(key)
	if !exist {
		return fallback
	}

	return value
}
