package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"gaetanjaminon/GoTuto/internal/catalog/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.CatalogConfig) (*gorm.DB, error) {
	// Get DSN from config with schema isolation
	dsn := cfg.Database.GetDSN()

	// Configure GORM logger based on config
	logLevel := logger.Info
	switch cfg.Logging.Level {
	case "debug":
		logLevel = logger.Info
	case "warn", "error":
		logLevel = logger.Warn
	default:
		logLevel = logger.Silent
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.WithContext(ctx).Exec("SELECT 1").Error; err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Store globally for easy access
	DB = db

	log.Printf("Catalog database connected successfully to %s:%d/%s (schema: %s)",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.Name, cfg.Database.Schema)
	return db, nil
}

// AutoMigrate runs GORM auto-migration for catalog models
// Note: For production, use the migration tool instead
func AutoMigrate(db *gorm.DB) error {
	// TODO: Add catalog models when they're created
	log.Println("Catalog database migration completed")
	return nil
}