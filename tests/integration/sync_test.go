package integration

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/harrybawsac/knmi-go/internal/cli"
	"github.com/harrybawsac/knmi-go/internal/db"
	_ "github.com/lib/pq"
)

// createMockKNMIServer creates a test server that serves KNMI-like data.
func createMockKNMIServer(t *testing.T) *httptest.Server {
	t.Helper()

	csvData := `# KNMI - Royal Netherlands Meteorological Institute
# Station 260 - De Bilt
# STN,YYYYMMDD,DDVEC,FHVEC,FG,FHX,FHXH,FHN,FHNH,FXX,FXXH,TG,TN,TNH,TX,TXH,T10N,T10NH,SQ,SP,Q,DR,RH,RHX,RHXH,PG,PX,PXH,PN,PNH,VVN,VVNH,VVX,VVXH,NG,UG,UX,UXH,UN,UNH,EV24
  260,20240101,  230,   45,   52,   72,   15,   31,    1,  100,   15,   85,   62,    6,  102,   14,   52,    6,   25,   28, 380,   10,   32,    8,   12,10250,10280,   12,10220,    6,   54,    7,   75,   15,    6,   88,   96,    7,   78,   14,    8
  260,20240102,  180,   38,   45,   65,   10,   25,    5,   90,   10,   90,   70,    3,  110,   15,   60,    3,   30,   35, 400,    5,   20,    5,    8,10260,10290,   10,10230,    5,   50,    5,   80,   12,    5,   85,   92,    5,   75,   12,   10
  260,20240103,  210,   42,   48,   68,   12,   28,    3,   95,   12,   88,   65,    5,  105,   14,   55,    5,   28,   32, 390,    8,   28,    7,   10,10255,10285,   11,10225,    6,   52,    6,   78,   14,    5,   86,   94,    6,   76,   13,    9
`

	// Create a zip file containing the CSV data
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	f, err := w.Create("etmgeg_260.txt")
	if err != nil {
		t.Fatalf("failed to create zip entry: %v", err)
	}
	if _, err := f.Write([]byte(csvData)); err != nil {
		t.Fatalf("failed to write to zip: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("failed to close zip: %v", err)
	}

	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/zip")
		rw.WriteHeader(http.StatusOK)
		rw.Write(buf.Bytes())
	}))
}

// applyMigrations runs the actual migrations.
func applyMigrations(t *testing.T, database *sql.DB, migrationsDir string) {
	t.Helper()

	// Create migrations table
	_, err := database.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			version INTEGER NOT NULL UNIQUE,
			filename VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("failed to create migrations table: %v", err)
	}

	// Create weather_records table
	_, err = database.Exec(`
		CREATE TABLE IF NOT EXISTS weather_records (
			id SERIAL PRIMARY KEY,
			station_id INTEGER NOT NULL,
			date DATE NOT NULL,
			ddvec INTEGER,
			fhvec INTEGER,
			fg INTEGER,
			fhx INTEGER,
			fhxh INTEGER,
			fhn INTEGER,
			fhnh INTEGER,
			fxx INTEGER,
			fxxh INTEGER,
			tg INTEGER,
			tn INTEGER,
			tnh INTEGER,
			tx INTEGER,
			txh INTEGER,
			t10n INTEGER,
			t10nh INTEGER,
			sq INTEGER,
			sp INTEGER,
			q INTEGER,
			dr INTEGER,
			rh INTEGER,
			rhx INTEGER,
			rhxh INTEGER,
			pg INTEGER,
			px INTEGER,
			pxh INTEGER,
			pn INTEGER,
			pnh INTEGER,
			vvn INTEGER,
			vvnh INTEGER,
			vvx INTEGER,
			vvxh INTEGER,
			ng INTEGER,
			ug INTEGER,
			ux INTEGER,
			uxh INTEGER,
			un INTEGER,
			unh INTEGER,
			ev24 INTEGER,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (station_id, date)
		)
	`)
	if err != nil {
		t.Fatalf("failed to create weather_records table: %v", err)
	}
}

func TestSyncCommand(t *testing.T) {
	databaseURL := getTestDatabaseURL(t)

	database, err := db.Connect(databaseURL)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	// Clean up before and after tests
	cleanupDatabase(t, database)
	t.Cleanup(func() { cleanupDatabase(t, database) })

	// Apply migrations
	applyMigrations(t, database, "")

	t.Run("syncs weather data successfully", func(t *testing.T) {
		server := createMockKNMIServer(t)
		defer server.Close()

		os.Setenv("DATABASE_URL", databaseURL)
		defer os.Unsetenv("DATABASE_URL")

		cmd := cli.NewRootCommand()
		cmd.SetArgs([]string{"sync", "--url", server.URL})

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("sync command failed: %v", err)
		}

		// Verify records were inserted
		var count int
		err = database.QueryRow("SELECT COUNT(*) FROM weather_records").Scan(&count)
		if err != nil {
			t.Fatalf("failed to count records: %v", err)
		}
		if count != 3 {
			t.Errorf("expected 3 records, got %d", count)
		}
	})

	t.Run("skips duplicate records on re-sync", func(t *testing.T) {
		server := createMockKNMIServer(t)
		defer server.Close()

		os.Setenv("DATABASE_URL", databaseURL)
		defer os.Unsetenv("DATABASE_URL")

		// Run sync twice
		cmd1 := cli.NewRootCommand()
		cmd1.SetArgs([]string{"sync", "--url", server.URL})
		if err := cmd1.Execute(); err != nil {
			t.Fatalf("first sync failed: %v", err)
		}

		cmd2 := cli.NewRootCommand()
		cmd2.SetArgs([]string{"sync", "--url", server.URL})
		if err := cmd2.Execute(); err != nil {
			t.Fatalf("second sync failed: %v", err)
		}

		// Verify no duplicates
		var count int
		err := database.QueryRow("SELECT COUNT(*) FROM weather_records").Scan(&count)
		if err != nil {
			t.Fatalf("failed to count records: %v", err)
		}
		if count != 3 {
			t.Errorf("expected 3 records (no duplicates), got %d", count)
		}
	})

	t.Run("reports error when table does not exist", func(t *testing.T) {
		// Drop the table to simulate missing migrations
		_, err := database.Exec("DROP TABLE IF EXISTS weather_records CASCADE")
		if err != nil {
			t.Fatalf("failed to drop table: %v", err)
		}

		server := createMockKNMIServer(t)
		defer server.Close()

		os.Setenv("DATABASE_URL", databaseURL)
		defer os.Unsetenv("DATABASE_URL")

		cmd := cli.NewRootCommand()
		cmd.SetArgs([]string{"sync", "--url", server.URL})

		err = cmd.Execute()
		if err == nil {
			t.Error("expected error when table doesn't exist, got nil")
		}
	})
}
