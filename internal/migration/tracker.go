// Package migration provides database migration functionality.
package migration

import (
	"database/sql"
	"fmt"
	"time"
)

// Tracker manages the state of applied migrations in the database.
type Tracker struct {
	db *sql.DB
}

// NewTracker creates a new migration tracker.
func NewTracker(db *sql.DB) *Tracker {
	return &Tracker{db: db}
}

// EnsureTable creates the migrations tracking table if it doesn't exist.
func (t *Tracker) EnsureTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			version INTEGER NOT NULL UNIQUE,
			filename VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := t.db.Exec(query)
	if err != nil {
		return fmt.Errorf("creating migrations table: %w", err)
	}
	return nil
}

// Applied returns a map of version numbers that have already been applied.
func (t *Tracker) Applied() (map[int]bool, error) {
	rows, err := t.db.Query("SELECT version FROM migrations")
	if err != nil {
		return nil, fmt.Errorf("querying applied migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[int]bool)
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("scanning migration version: %w", err)
		}
		applied[version] = true
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating migrations: %w", err)
	}

	return applied, nil
}

// Record marks a migration as applied.
func (t *Tracker) Record(version int, filename string) error {
	query := `INSERT INTO migrations (version, filename, applied_at) VALUES ($1, $2, $3)`
	_, err := t.db.Exec(query, version, filename, time.Now())
	if err != nil {
		return fmt.Errorf("recording migration: %w", err)
	}
	return nil
}
