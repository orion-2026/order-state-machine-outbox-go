package repository

import (
	"fmt"
	"order-state-machine-outbox-go/internal/config"
	"order-state-machine-outbox-go/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Open(cfg config.Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	switch cfg.DatabaseProvider {
	case "postgres":
		db, err = gorm.Open(postgres.Open(cfg.DatabaseConnectionString), &gorm.Config{})
	default:
		db, err = gorm.Open(sqlite.Open(cfg.DatabaseConnectionString), &gorm.Config{})
	}
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.AutoMigrate(&domain.Order{}, &domain.OutboxEvent{}); err != nil {
		return nil, fmt.Errorf("migrate database: %w", err)
	}

	return db, nil
}
