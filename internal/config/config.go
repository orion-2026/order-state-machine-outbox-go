package config

import "os"

type Config struct {
	Port                     string
	DatabaseProvider         string
	DatabaseConnectionString string
}

func Load() Config {
	port := getenv("PORT", "8080")
	provider := getenv("DATABASE_PROVIDER", "sqlite")
	connectionString := getenv("DATABASE_CONNECTION_STRING", "order-demo.db")
	return Config{
		Port:                     port,
		DatabaseProvider:         provider,
		DatabaseConnectionString: connectionString,
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
