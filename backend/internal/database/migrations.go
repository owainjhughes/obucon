package database

import (
	"fmt"
	"log"
	"obucon/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// TODO: Implement RunMigrations function
//
// This function should:
// 1. Create a migration instance pointing to your migrations directory
// 2. Run any pending migrations up to the latest version
// 3. Handle errors gracefully (distinguish between "no migrations" vs actual errors)
//
// Implementation notes:
//   - Migration source should be: file:///path/to/backend/migrations
//   - Database DSN format: postgres://user:password@host:port/dbname?sslmode=disable
//   - Use migrate.Up() to run all pending migrations
//   - Check for migrate.ErrNoChange (not an error - means already up to date)
//
// Reference:
//   - golang-migrate docs: https://github.com/golang-migrate/migrate
//   - Examples: https://github.com/golang-migrate/migrate/tree/master/example
//
// Hint: Wrap the migration instance in error handling:
//

func RunMigrations(cfg *config.Config) error {
	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	// Use relative path for migrations
	// golang-migrate on Windows has issues with absolute paths
	sourceURL := "file://migrations"

	log.Printf("Migration source: %s", sourceURL)
	log.Printf("Database: %s:%s/%s", cfg.DBHost, cfg.DBPort, cfg.DBName)

	// Create migration instance
	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	defer m.Close()

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}
