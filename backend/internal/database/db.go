package database

import (
	"fmt"
	"log"
	"obucon/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB initializes the database connection
//
// This function:
// 1. Builds PostgreSQL connection string from config
// 2. Connects to the database using GORM
// 3. Tests the connection with a ping
// 4. Returns a GORM DB instance
//
// TODO: Add migration runner to automatically execute .sql files
// Consider using github.com/golang-migrate/migrate for migration management
//
// Reference:
// - GORM PostgreSQL docs: https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL
// - golang-migrate: https://github.com/golang-migrate/migrate
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName,
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
