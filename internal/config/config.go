package config

import "os"

type Config struct {
	Port           string
	DatabaseDriver string
	DatabaseDSN    string
}

func Load() Config {
	port := getenv("PORT", "8080")
	driver := getenv("DATABASE_DRIVER", "sqlite")
	dsn := getenv("DATABASE_DSN", "order-demo.db")
	return Config{
		Port:           port,
		DatabaseDriver: driver,
		DatabaseDSN:    dsn,
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
