package cli

import (
	"bytes"
	"fmt"

	"github.com/harrybawsac/knmi-go/internal/db"
	"github.com/harrybawsac/knmi-go/internal/fetch"
	"github.com/harrybawsac/knmi-go/internal/parser"
	"github.com/spf13/cobra"
)

var dataURL string
var dryRun bool

// newSyncCommand creates the sync subcommand.
func newSyncCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Download and sync KNMI weather data",
		Long: `Download weather data from KNMI and sync to the database.

The command downloads a zip file from the KNMI website, extracts the CSV data,
and inserts new records into the database. Existing records are skipped.`,
		RunE: runSync,
	}

	cmd.Flags().StringVar(&dataURL, "url", "", "Override KNMI data URL")
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "Preview records without inserting")

	return cmd
}

// formatPreviewValue formats an optional int value for display.
// Returns "-" for nil values, otherwise the integer as a string.
func formatPreviewValue(v *int) string {
	if v == nil {
		return "-"
	}
	return fmt.Sprintf("%d", *v)
}

// printPreviewTable prints a tabular preview of weather records.
// Handles empty result case per FR-008.
func printPreviewTable(records []parser.WeatherRecord, total int) {
	if len(records) == 0 {
		fmt.Println("Dry-run mode: no new records to insert")
		fmt.Println()
		fmt.Println("All records from the KNMI file already exist in the database.")
		return
	}

	fmt.Println("Dry-run mode: previewing records that would be inserted")
	fmt.Println()

	// Print header
	fmt.Printf("%-10s %10s %6s %6s %6s %6s %6s\n", "DATE", "STATION_ID", "TG", "TN", "TX", "FG", "RH")

	// Get last 10 (or fewer) records
	start := 0
	if len(records) > 10 {
		start = len(records) - 10
	}
	previewRecords := records[start:]

	// Print records
	for _, rec := range previewRecords {
		fmt.Printf("%-10s %10d %6s %6s %6s %6s %6s\n",
			rec.Date.Format("2006-01-02"),
			rec.StationID,
			formatPreviewValue(rec.TG),
			formatPreviewValue(rec.TN),
			formatPreviewValue(rec.TX),
			formatPreviewValue(rec.FG),
			formatPreviewValue(rec.RH),
		)
	}

	fmt.Println()
	if total <= 10 {
		fmt.Printf("Total: %d new records would be inserted (showing all)\n", total)
	} else {
		fmt.Printf("Total: %d new records would be inserted (showing last 10)\n", total)
	}
}

// runSync executes the sync command.
func runSync(cmd *cobra.Command, args []string) error {
	cfg := GetConfig()
	dbURL := cfg.DatabaseURL
	if databaseURL != "" {
		dbURL = databaseURL
	}

	// Determine data URL
	url := dataURL
	if url == "" {
		url = cfg.KNMIDataURL
	}

	// Download data
	LogVerbose("Downloading from %s...", url)
	zipData, err := fetch.Download(url)
	if err != nil {
		return fmt.Errorf("failed to download data: %w", err)
	}
	LogVerbose("Downloaded %.2f MB", float64(len(zipData))/(1024*1024))

	// Extract zip
	LogVerbose("Extracting archive...")
	csvData, err := fetch.ExtractZip(zipData)
	if err != nil {
		return fmt.Errorf("failed to extract zip: %w", err)
	}

	// Parse CSV
	LogVerbose("Parsing CSV...")
	records, err := parser.ParseCSV(bytes.NewReader(csvData))
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %w", err)
	}
	LogVerbose("Parsed %d rows", len(records))

	// Dry-run mode: preview without inserting
	if dryRun {
		// If database is configured, filter to show only new records
		if dbURL != "" {
			LogVerbose("Connecting to database for duplicate filtering...")
			database, err := db.Connect(dbURL)
			if err != nil {
				LogVerbose("Warning: could not connect to database, showing all parsed records")
				printPreviewTable(records, len(records))
				return nil
			}
			defer database.Close()

			repo := db.NewWeatherRepository(database)
			tableExists, err := repo.TableExists()
			if err != nil || !tableExists {
				LogVerbose("Warning: table not found, showing all parsed records")
				printPreviewTable(records, len(records))
				return nil
			}

			LogVerbose("Dry-run mode: filtering new records...")
			newRecords, err := repo.FilterNewRecords(records)
			if err != nil {
				return fmt.Errorf("filtering new records: %w", err)
			}
			printPreviewTable(newRecords, len(newRecords))
			return nil
		}

		// No database configured - just show parsed records
		LogVerbose("Dry-run mode: no database configured, showing parsed records")
		printPreviewTable(records, len(records))
		return nil
	}

	// Normal sync mode - database is required
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

	// Check if migrations have been applied
	repo := db.NewWeatherRepository(database)
	tableExists, err := repo.TableExists()
	if err != nil {
		return fmt.Errorf("checking database state: %w", err)
	}
	if !tableExists {
		return fmt.Errorf("no migrations applied. Run 'knmi migrate' first")
	}

	// Insert records
	LogVerbose("Inserting records...")
	result, err := repo.InsertRecords(records)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	// Get total count
	total, err := repo.GetTotalCount()
	if err != nil {
		LogVerbose("Warning: could not get total count: %v", err)
		total = result.Inserted
	}

	// Print summary
	fmt.Printf("Synced %d new records (%d total)\n", result.Inserted, total)

	return nil
}
