package config

import (
	"errors"
	"fmt"
	"os"
	"time"
)

// Config holds the settings read from the environment. Only cmd/api/main.go
// imports this package: adapters receive concrete values through constructors.
type Config struct {
	HTTPAddr  string
	MySQLDSN  string
	JWTSecret string
	JWTTTL    time.Duration
}

// Load reads the configuration from environment variables.
func Load() (Config, error) {
	mysqlDSN := os.Getenv("MYSQL_DSN")
	jwtSecret := os.Getenv("JWT_SECRET")
	if mysqlDSN == "" || jwtSecret == "" {
		return Config{}, errors.New("MYSQL_DSN and JWT_SECRET are required")
	}

	jwtTTL, err := time.ParseDuration(envOrDefault("JWT_TTL", "15m"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid JWT_TTL: %w", err)
	}

	return Config{
		HTTPAddr:  envOrDefault("HTTP_ADDR", ":8080"),
		MySQLDSN:  mysqlDSN,
		JWTSecret: jwtSecret,
		JWTTTL:    jwtTTL,
	}, nil
}

func envOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
