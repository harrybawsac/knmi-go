package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/harrybawsac/knmi-go/internal/parser"
)

// WeatherRepository manages weather records in the database.
type WeatherRepository struct {
	db *sql.DB
}

// NewWeatherRepository creates a new weather repository.
func NewWeatherRepository(db *sql.DB) *WeatherRepository {
	return &WeatherRepository{db: db}
}

// InsertResult contains the result of an insert operation.
type InsertResult struct {
	Inserted int
	Skipped  int
	Total    int
}

// InsertRecords inserts weather records into the database.
// Uses ON CONFLICT DO NOTHING to skip duplicates.
func (r *WeatherRepository) InsertRecords(records []parser.WeatherRecord) (*InsertResult, error) {
	result := &InsertResult{
		Total: len(records),
	}

	if len(records) == 0 {
		return result, nil
	}

	query := `
		INSERT INTO weather_records (
			station_id, date, ddvec, fhvec, fg, fhx, fhxh, fhn, fhnh, fxx, fxxh,
			tg, tn, tnh, tx, txh, t10n, t10nh, sq, sp, q,
			dr, rh, rhx, rhxh, pg, px, pxh, pn, pnh,
			vvn, vvnh, vvx, vvxh, ng, ug, ux, uxh, un, unh, ev24
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11,
			$12, $13, $14, $15, $16, $17, $18, $19, $20, $21,
			$22, $23, $24, $25, $26, $27, $28, $29, $30,
			$31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41
		)
		ON CONFLICT (station_id, date) DO NOTHING
	`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("preparing insert statement: %w", err)
	}
	defer stmt.Close()

	for _, rec := range records {
		res, err := stmt.Exec(
			rec.StationID, rec.Date,
			rec.DDVEC, rec.FHVEC, rec.FG, rec.FHX, rec.FHXH, rec.FHN, rec.FHNH, rec.FXX, rec.FXXH,
			rec.TG, rec.TN, rec.TNH, rec.TX, rec.TXH, rec.T10N, rec.T10NH, rec.SQ, rec.SP, rec.Q,
			rec.DR, rec.RH, rec.RHX, rec.RHXH, rec.PG, rec.PX, rec.PXH, rec.PN, rec.PNH,
			rec.VVN, rec.VVNH, rec.VVX, rec.VVXH, rec.NG, rec.UG, rec.UX, rec.UXH, rec.UN, rec.UNH, rec.EV24,
		)
		if err != nil {
			return result, fmt.Errorf("inserting record for date %s: %w", rec.Date.Format("2006-01-02"), err)
		}

		affected, err := res.RowsAffected()
		if err != nil {
			return result, fmt.Errorf("getting rows affected: %w", err)
		}

		if affected > 0 {
			result.Inserted++
		} else {
			result.Skipped++
		}
	}

	return result, nil
}

// GetTotalCount returns the total number of weather records.
func (r *WeatherRepository) GetTotalCount() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM weather_records").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting weather records: %w", err)
	}
	return count, nil
}

// GetLatestDate returns the most recent date in the database.
func (r *WeatherRepository) GetLatestDate() (*time.Time, error) {
	var date sql.NullTime
	err := r.db.QueryRow("SELECT MAX(date) FROM weather_records").Scan(&date)
	if err != nil {
		return nil, fmt.Errorf("getting latest date: %w", err)
	}
	if !date.Valid {
		return nil, nil
	}
	return &date.Time, nil
}

// TableExists checks if the weather_records table exists.
func (r *WeatherRepository) TableExists() (bool, error) {
	var exists bool
	err := r.db.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_name = 'weather_records'
		)
	`).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking weather_records table: %w", err)
	}
	return exists, nil
}

// FilterNewRecords returns only records that don't already exist in the database.
// Uses the (station_id, date) composite key to check for existing records.
func (r *WeatherRepository) FilterNewRecords(records []parser.WeatherRecord) ([]parser.WeatherRecord, error) {
	if len(records) == 0 {
		return records, nil
	}

	// Query all existing (station_id, date) pairs
	rows, err := r.db.Query("SELECT station_id, date FROM weather_records")
	if err != nil {
		return nil, fmt.Errorf("querying existing records: %w", err)
	}
	defer rows.Close()

	// Build a set of existing keys
	existing := make(map[string]struct{})
	for rows.Next() {
		var stationID int
		var date time.Time
		if err := rows.Scan(&stationID, &date); err != nil {
			return nil, fmt.Errorf("scanning existing record: %w", err)
		}
		key := fmt.Sprintf("%d:%s", stationID, date.Format("2006-01-02"))
		existing[key] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating existing records: %w", err)
	}

	// Filter to only new records
	var newRecords []parser.WeatherRecord
	for _, rec := range records {
		key := fmt.Sprintf("%d:%s", rec.StationID, rec.Date.Format("2006-01-02"))
		if _, exists := existing[key]; !exists {
			newRecords = append(newRecords, rec)
		}
	}

	return newRecords, nil
}
