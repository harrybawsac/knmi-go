package unit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/harrybawsac/knmi-go/internal/migration"
)

func TestDiscoverMigrations(t *testing.T) {
	// Create a temporary directory with test migration files
	tempDir := t.TempDir()

	testCases := []struct {
		name          string
		files         []string
		expectedCount int
		expectedOrder []int
		expectError   bool
	}{
		{
			name:          "no migrations",
			files:         []string{},
			expectedCount: 0,
			expectedOrder: []int{},
			expectError:   false,
		},
		{
			name:          "single migration",
			files:         []string{"001_create_tables.sql"},
			expectedCount: 1,
			expectedOrder: []int{1},
			expectError:   false,
		},
		{
			name:          "multiple migrations in order",
			files:         []string{"001_create_tables.sql", "002_add_indexes.sql", "003_add_constraints.sql"},
			expectedCount: 3,
			expectedOrder: []int{1, 2, 3},
			expectError:   false,
		},
		{
			name:          "migrations out of order should be sorted",
			files:         []string{"003_third.sql", "001_first.sql", "002_second.sql"},
			expectedCount: 3,
			expectedOrder: []int{1, 2, 3},
			expectError:   false,
		},
		{
			name:          "ignore non-sql files",
			files:         []string{"001_create_tables.sql", "README.md", "002_indexes.sql"},
			expectedCount: 2,
			expectedOrder: []int{1, 2},
			expectError:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a subdirectory for this test case
			testDir := filepath.Join(tempDir, tc.name)
			if err := os.MkdirAll(testDir, 0755); err != nil {
				t.Fatalf("failed to create test directory: %v", err)
			}

			// Create test files
			for _, f := range tc.files {
				path := filepath.Join(testDir, f)
				if err := os.WriteFile(path, []byte("-- test"), 0644); err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
			}

			// Run discovery
			migrations, err := migration.Discover(testDir)

			// Check error expectation
			if tc.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Check count
			if len(migrations) != tc.expectedCount {
				t.Errorf("expected %d migrations, got %d", tc.expectedCount, len(migrations))
			}

			// Check order
			for i, expectedVersion := range tc.expectedOrder {
				if i >= len(migrations) {
					break
				}
				if migrations[i].Version != expectedVersion {
					t.Errorf("migration %d: expected version %d, got %d", i, expectedVersion, migrations[i].Version)
				}
			}
		})
	}
}

func TestParseVersion(t *testing.T) {
	testCases := []struct {
		filename        string
		expectedVersion int
		expectError     bool
	}{
		{"001_create_tables.sql", 1, false},
		{"002_add_indexes.sql", 2, false},
		{"010_something.sql", 10, false},
		{"100_big_migration.sql", 100, false},
		{"invalid.sql", 0, true},
		{"abc_not_a_number.sql", 0, true},
		{"", 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			version, err := migration.ParseVersion(tc.filename)

			if tc.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if version != tc.expectedVersion {
				t.Errorf("expected version %d, got %d", tc.expectedVersion, version)
			}
		})
	}
}
