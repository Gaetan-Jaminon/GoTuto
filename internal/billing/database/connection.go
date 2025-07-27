package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"gaetanjaminon/GoTuto/internal/billing/config"
	"gaetanjaminon/GoTuto/internal/billing/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.BillingConfig) (*gorm.DB, error) {
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

	// Test connection with configurable timeout
	timeout := cfg.Database.ConnectionTimeout
	if timeout == 0 {
		timeout = 5 * time.Second // Default fallback
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := db.WithContext(ctx).Exec("SELECT 1").Error; err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Billing database connected successfully to %s:%d/%s (schema: %s)",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.Name, cfg.Database.Schema)
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.Client{},
		&models.Invoice{},
	)

	if err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	log.Println("Database migration completed")
	return nil
}