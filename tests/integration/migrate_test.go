package integration

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/harrybawsac/knmi-go/internal/cli"
	"github.com/harrybawsac/knmi-go/internal/db"
	_ "github.com/lib/pq"
)

// getTestDatabaseURL returns the test database URL from environment or skips the test.
func getTestDatabaseURL(t *testing.T) string {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping integration test")
	}
	return url
}

// setupTestMigrations creates a temporary directory with test migration files.
func setupTestMigrations(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	// Create a simple migration file
	migration := `-- Create test table
CREATE TABLE IF NOT EXISTS test_table (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);`

	err := os.WriteFile(filepath.Join(dir, "001_create_test.sql"), []byte(migration), 0644)
	if err != nil {
		t.Fatalf("failed to write test migration: %v", err)
	}

	return dir
}

// cleanupDatabase drops test tables created during tests.
func cleanupDatabase(t *testing.T, database *sql.DB) {
	t.Helper()
	tables := []string{"test_table", "weather_records", "migrations"}
	for _, table := range tables {
		_, err := database.Exec("DROP TABLE IF EXISTS " + table + " CASCADE")
		if err != nil {
			t.Logf("warning: failed to drop table %s: %v", table, err)
		}
	}
}

func TestMigrateCommand(t *testing.T) {
	databaseURL := getTestDatabaseURL(t)

	database, err := db.Connect(databaseURL)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	// Clean up before and after tests
	cleanupDatabase(t, database)
	t.Cleanup(func() { cleanupDatabase(t, database) })

	t.Run("applies migrations successfully", func(t *testing.T) {
		migrationsDir := setupTestMigrations(t)

		// Set up environment
		os.Setenv("DATABASE_URL", databaseURL)
		os.Setenv("KNMI_MIGRATIONS_DIR", migrationsDir)
		defer os.Unsetenv("DATABASE_URL")
		defer os.Unsetenv("KNMI_MIGRATIONS_DIR")

		// Execute the migrate command
		cmd := cli.NewRootCommand()
		cmd.SetArgs([]string{"migrate", "--migrations-dir", migrationsDir})

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("migrate command failed: %v", err)
		}

		// Verify the test table was created
		var exists bool
		err = database.QueryRow(`
			SELECT EXISTS (
				SELECT FROM information_schema.tables 
				WHERE table_name = 'test_table'
			)
		`).Scan(&exists)
		if err != nil {
			t.Fatalf("failed to check if table exists: %v", err)
		}
		if !exists {
			t.Error("expected test_table to exist after migration")
		}

		// Verify migration was tracked
		var count int
		err = database.QueryRow("SELECT COUNT(*) FROM migrations WHERE filename = '001_create_test.sql'").Scan(&count)
		if err != nil {
			t.Fatalf("failed to check migrations table: %v", err)
		}
		if count != 1 {
			t.Errorf("expected migration to be tracked, got count=%d", count)
		}
	})

	t.Run("skips already applied migrations", func(t *testing.T) {
		migrationsDir := setupTestMigrations(t)

		os.Setenv("DATABASE_URL", databaseURL)
		os.Setenv("KNMI_MIGRATIONS_DIR", migrationsDir)
		defer os.Unsetenv("DATABASE_URL")
		defer os.Unsetenv("KNMI_MIGRATIONS_DIR")

		// Run migrate twice
		cmd1 := cli.NewRootCommand()
		cmd1.SetArgs([]string{"migrate", "--migrations-dir", migrationsDir})
		if err := cmd1.Execute(); err != nil {
			t.Fatalf("first migrate failed: %v", err)
		}

		cmd2 := cli.NewRootCommand()
		cmd2.SetArgs([]string{"migrate", "--migrations-dir", migrationsDir})
		if err := cmd2.Execute(); err != nil {
			t.Fatalf("second migrate failed: %v", err)
		}

		// Verify migration was only tracked once
		var count int
		err := database.QueryRow("SELECT COUNT(*) FROM migrations WHERE filename = '001_create_test.sql'").Scan(&count)
		if err != nil {
			t.Fatalf("failed to check migrations table: %v", err)
		}
		if count != 1 {
			t.Errorf("expected migration to be tracked once, got count=%d", count)
		}
	})

	t.Run("reports error for invalid SQL", func(t *testing.T) {
		dir := t.TempDir()
		invalidSQL := "THIS IS NOT VALID SQL !!!"
		err := os.WriteFile(filepath.Join(dir, "001_invalid.sql"), []byte(invalidSQL), 0644)
		if err != nil {
			t.Fatalf("failed to write invalid migration: %v", err)
		}

		os.Setenv("DATABASE_URL", databaseURL)
		defer os.Unsetenv("DATABASE_URL")

		cmd := cli.NewRootCommand()
		cmd.SetArgs([]string{"migrate", "--migrations-dir", dir})

		err = cmd.Execute()
		if err == nil {
			t.Error("expected error for invalid SQL, got nil")
		}
	})

	t.Run("reports error for missing migrations directory", func(t *testing.T) {
		os.Setenv("DATABASE_URL", databaseURL)
		defer os.Unsetenv("DATABASE_URL")

		cmd := cli.NewRootCommand()
		cmd.SetArgs([]string{"migrate", "--migrations-dir", "/nonexistent/path"})

		err := cmd.Execute()
		if err == nil {
			t.Error("expected error for missing directory, got nil")
		}
		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("expected error to mention 'not found', got: %v", err)
		}
	})
}
