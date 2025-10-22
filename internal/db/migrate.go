package db

import (
	"database/sql"
	"fmt"
	"sort"

	_ "embed"
)

//go:embed migrate/0001_init.sql
var migration0001 string

//go:embed migrate/0002_add_oauth_columns.sql
var migration0002 string

//go:embed migrate/0003_add_username.sql
var migration0003 string

// RunMigrations runs all database migrations
func RunMigrations(db *sql.DB) error {
	// Create migrations table if it doesn't exist
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	migrations := map[string]string{
		"0001_init.sql":              migration0001,
		"0002_add_oauth_columns.sql": migration0002,
		"0003_add_username.sql":      migration0003,
	}

	// Get list of migration files and sort them
	var migrationFiles []string
	for filename := range migrations {
		migrationFiles = append(migrationFiles, filename)
	}
	sort.Strings(migrationFiles)

	for _, filename := range migrationFiles {
		// Check if migration has already been applied
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = $1", filename).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check migration status for %s: %w", filename, err)
		}

		if count > 0 {
			fmt.Printf("Migration %s already applied, skipping\n", filename)
			continue
		}

		// Apply migration
		fmt.Printf("Applying migration %s\n", filename)
		_, err = db.Exec(migrations[filename])
		if err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", filename, err)
		}

		// Record migration as applied
		_, err = db.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", filename)
		if err != nil {
			return fmt.Errorf("failed to record migration %s: %w", filename, err)
		}

		fmt.Printf("Migration %s applied successfully\n", filename)
	}

	return nil
}