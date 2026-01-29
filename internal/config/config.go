// Package config provides configuration management for the KNMI CLI.
package config

import (
	"os"
)

const (
	// DefaultKNMIDataURL is the default URL for KNMI weather data.
	DefaultKNMIDataURL = "https://cdn.knmi.nl/knmi/map/page/klimatologie/gegevens/daggegevens/etmgeg_260.zip"

	// DefaultMigrationsDir is the default directory for SQL migration files.
	DefaultMigrationsDir = "./migrations"
)

// Config holds the application configuration.
type Config struct {
	// DatabaseURL is the PostgreSQL connection string.
	DatabaseURL string

	// KNMIDataURL is the URL to fetch KNMI weather data from.
	KNMIDataURL string

	// MigrationsDir is the path to the migrations directory.
	MigrationsDir string

	// Verbose enables detailed logging output.
	Verbose bool
}

// Load creates a Config from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		DatabaseURL:   getEnv("DATABASE_URL", ""),
		KNMIDataURL:   getEnv("KNMI_DATA_URL", DefaultKNMIDataURL),
		MigrationsDir: getEnv("KNMI_MIGRATIONS_DIR", DefaultMigrationsDir),
		Verbose:       true,
	}
}

// getEnv returns the value of an environment variable or a default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
