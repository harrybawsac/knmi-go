// Package migration provides database migration functionality.
package migration

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Migration represents a single database migration file.
type Migration struct {
	Version  int
	Name     string
	Filename string
	Path     string
	Content  string
}

// versionRegex matches migration filenames like "001_create_tables.sql"
var versionRegex = regexp.MustCompile(`^(\d+)_.*\.sql$`)

// Discover finds all migration files in the given directory and returns them sorted by version.
func Discover(dir string) ([]Migration, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("migrations directory not found: %s", dir)
		}
		return nil, fmt.Errorf("reading migrations directory: %w", err)
	}

	var migrations []Migration
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if !strings.HasSuffix(filename, ".sql") {
			continue
		}

		version, err := ParseVersion(filename)
		if err != nil {
			// Skip files that don't match the expected pattern
			continue
		}

		path := filepath.Join(dir, filename)
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading migration file %s: %w", filename, err)
		}

		// Extract name from filename (e.g., "001_create_tables.sql" -> "create_tables")
		name := strings.TrimSuffix(filename, ".sql")
		if idx := strings.Index(name, "_"); idx >= 0 {
			name = name[idx+1:]
		}

		migrations = append(migrations, Migration{
			Version:  version,
			Name:     name,
			Filename: filename,
			Path:     path,
			Content:  string(content),
		})
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// ParseVersion extracts the version number from a migration filename.
func ParseVersion(filename string) (int, error) {
	if filename == "" {
		return 0, fmt.Errorf("empty filename")
	}

	matches := versionRegex.FindStringSubmatch(filename)
	if matches == nil {
		return 0, fmt.Errorf("invalid migration filename format: %s", filename)
	}

	version, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("parsing version number: %w", err)
	}

	return version, nil
}

// Result represents the result of running migrations.
type Result struct {
	Applied  []string
	Skipped  int
	Duration time.Duration
}

// LogFunc is a function type for logging messages.
type LogFunc func(format string, args ...interface{})

// Runner executes database migrations.
type Runner struct {
	db      *sql.DB
	tracker *Tracker
	logFn   LogFunc
}

// NewRunner creates a new migration runner.
func NewRunner(db *sql.DB, logFn LogFunc) *Runner {
	return &Runner{
		db:      db,
		tracker: NewTracker(db),
		logFn:   logFn,
	}
}

// Run discovers and applies pending migrations from the given directory.
func (r *Runner) Run(dir string) (*Result, error) {
	start := time.Now()

	// Ensure migrations table exists
	if err := r.tracker.EnsureTable(); err != nil {
		return nil, fmt.Errorf("ensuring migrations table: %w", err)
	}

	// Discover migrations
	r.log("Discovering migrations in %s...", dir)
	migrations, err := Discover(dir)
	if err != nil {
		return nil, err
	}

	if len(migrations) == 0 {
		r.log("No migrations found")
		return &Result{Duration: time.Since(start)}, nil
	}

	// Get already applied migrations
	applied, err := r.tracker.Applied()
	if err != nil {
		return nil, fmt.Errorf("getting applied migrations: %w", err)
	}

	// Filter to pending migrations
	var pending []Migration
	for _, m := range migrations {
		if !applied[m.Version] {
			pending = append(pending, m)
		}
	}

	if len(pending) == 0 {
		r.log("All migrations already applied")
		return &Result{
			Skipped:  len(migrations),
			Duration: time.Since(start),
		}, nil
	}

	r.log("Found %d pending migrations", len(pending))

	// Apply each pending migration
	result := &Result{}
	for _, m := range pending {
		if err := r.applyMigration(m); err != nil {
			return result, fmt.Errorf("migration %s failed: %w", m.Filename, err)
		}
		result.Applied = append(result.Applied, m.Filename)
	}

	result.Duration = time.Since(start)
	result.Skipped = len(migrations) - len(pending)

	return result, nil
}

// applyMigration applies a single migration within a transaction.
func (r *Runner) applyMigration(m Migration) error {
	r.log("Applying %s...", m.Filename)
	migrationStart := time.Now()

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback()

	// Execute the migration SQL
	if _, err := tx.Exec(m.Content); err != nil {
		return fmt.Errorf("executing SQL: %w", err)
	}

	// Record the migration
	if err := r.tracker.Record(m.Version, m.Filename); err != nil {
		return fmt.Errorf("recording migration: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	r.log("Applied %s (%v)", m.Filename, time.Since(migrationStart).Round(time.Millisecond))
	return nil
}

// log prints a message if a log function is configured.
func (r *Runner) log(format string, args ...interface{}) {
	if r.logFn != nil {
		r.logFn(format, args...)
	}
}
