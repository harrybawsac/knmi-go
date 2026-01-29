// Package cli provides the command-line interface for the KNMI tool.
package cli

import (
	"fmt"
	"os"

	"github.com/harrybawsac/knmi-go/internal/config"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// Version is set at build time.
var Version = "dev"

var (
	cfg         *config.Config
	verbose     bool
	databaseURL string
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "knmi",
	Short: "KNMI weather data sync tool",
	Long: `A CLI tool to fetch weather data from KNMI (Royal Netherlands 
Meteorological Institute) and sync it to a PostgreSQL database.

The tool supports:
  - Database migrations for schema management
  - Incremental data sync (only new records are inserted)
  - Configurable data source URL

Environment Variables:
  DATABASE_URL         PostgreSQL connection string
  KNMI_DATA_URL        Override default KNMI data URL
  KNMI_MIGRATIONS_DIR  Path to migrations directory`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// NewRootCommand creates and returns a new root command for testing.
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "knmi",
		Short: "KNMI weather data sync tool",
		Long: `A CLI tool to fetch weather data from KNMI (Royal Netherlands 
Meteorological Institute) and sync it to a PostgreSQL database.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "", false, "Enable detailed progress logging")
	cmd.PersistentFlags().StringVar(&databaseURL, "database-url", "", "PostgreSQL connection string (overrides DATABASE_URL)")

	// Add subcommands
	cmd.AddCommand(newMigrateCommand())
	cmd.AddCommand(newSyncCommand())
	cmd.AddCommand(newVersionCommand())

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "", false, "Enable detailed progress logging")
	rootCmd.PersistentFlags().StringVar(&databaseURL, "database-url", "", "PostgreSQL connection string (overrides DATABASE_URL)")

	// Add subcommands
	rootCmd.AddCommand(newMigrateCommand())
	rootCmd.AddCommand(newSyncCommand())
	rootCmd.AddCommand(newVersionCommand())
}

// newVersionCommand creates the version subcommand.
func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Long:  `Display the version of the knmi CLI tool.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("knmi version %s\n", Version)
		},
	}
}

// initConfig reads in config from environment variables.
func initConfig() {
	// Load .env file if it exists (silently ignore if missing)
	_ = godotenv.Load()

	cfg = config.Load()
	cfg.Verbose = verbose

	// Command-line flag overrides environment variable
	if databaseURL != "" {
		cfg.DatabaseURL = databaseURL
	}
}

// GetConfig returns the current configuration.
func GetConfig() *config.Config {
	return cfg
}

// IsVerbose returns whether verbose mode is enabled.
func IsVerbose() bool {
	return verbose
}

// LogVerbose prints a message only if verbose mode is enabled.
func LogVerbose(format string, args ...interface{}) {
	if verbose {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}
}
