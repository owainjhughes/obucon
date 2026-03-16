package database

import (
	"fmt"
	"log"
	"obucon/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// - GORM PostgreSQL docs: https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL
// - golang-migrate: https://github.com/golang-migrate/migrate
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBSSLMode,
	)

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying database connection: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Database connection established successfully to %s:%s/%s", cfg.DBHost, cfg.DBPort, cfg.DBName)

	return db, nil
}
