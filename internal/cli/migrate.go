package cli

import (
	"fmt"

	"github.com/harrybawsac/knmi-go/internal/db"
	"github.com/harrybawsac/knmi-go/internal/migration"
	"github.com/spf13/cobra"
)

var migrationsDir string

// newMigrateCommand creates the migrate subcommand.
func newMigrateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Apply pending database migrations",
		Long: `Apply pending database migrations from the migrations directory.

Migrations are SQL files named with a version prefix (e.g., 001_create_tables.sql).
Each migration is applied in a transaction and tracked in the 'migrations' table.`,
		RunE: runMigrate,
	}

	cmd.Flags().StringVar(&migrationsDir, "migrations-dir", "./migrations", "Path to migrations directory")

	return cmd
}

// runMigrate executes the migrate command.
func runMigrate(cmd *cobra.Command, args []string) error {
	// Get database URL from config or flag
	cfg := GetConfig()
	dbURL := cfg.DatabaseURL
	if databaseURL != "" {
		dbURL = databaseURL
	}

	if dbURL == "" {
		return fmt.Errorf("database URL not configured (set DATABASE_URL or use --database-url)")
	}

	// Connect to database
	LogVerbose("Connecting to database...")
	database, err := db.Connect(dbURL)
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}
	defer database.Close()

	// Use migrations directory from flag or config
	dir := migrationsDir
	if dir == "./migrations" && cfg.MigrationsDir != "" {
		dir = cfg.MigrationsDir
	}

	// Create and run migrations
	var logFn migration.LogFunc
	if IsVerbose() {
		logFn = func(format string, args ...interface{}) {
			fmt.Printf(format+"\n", args...)
		}
	}

	runner := migration.NewRunner(database, logFn)
	result, err := runner.Run(dir)
	if err != nil {
		return err
	}

	// Print summary
	if len(result.Applied) == 0 {
		fmt.Println("No migrations to apply")
	} else {
		fmt.Printf("Applied %d migrations:\n", len(result.Applied))
		for _, name := range result.Applied {
			fmt.Printf("  %s\n", name)
		}
	}

	return nil
}
