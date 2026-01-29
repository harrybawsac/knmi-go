// Package parser provides CSV parsing for KNMI weather data.
package parser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// ExpectedColumns is the number of columns in the KNMI data file.
const ExpectedColumns = 41

// WeatherRecord represents a single weather observation.
type WeatherRecord struct {
	StationID int
	Date      time.Time
	DDVEC     *int
	FHVEC     *int
	FG        *int
	FHX       *int
	FHXH      *int
	FHN       *int
	FHNH      *int
	FXX       *int
	FXXH      *int
	TG        *int
	TN        *int
	TNH       *int
	TX        *int
	TXH       *int
	T10N      *int
	T10NH     *int
	SQ        *int
	SP        *int
	Q         *int
	DR        *int
	RH        *int
	RHX       *int
	RHXH      *int
	PG        *int
	PX        *int
	PXH       *int
	PN        *int
	PNH       *int
	VVN       *int
	VVNH      *int
	VVX       *int
	VVXH      *int
	NG        *int
	UG        *int
	UX        *int
	UXH       *int
	UN        *int
	UNH       *int
	EV24      *int
}

// ParseCSV parses KNMI weather data from a reader.
func ParseCSV(r io.Reader) ([]WeatherRecord, error) {
	var records []WeatherRecord

	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Skip header/description lines - data lines start with station ID (digits)
		// KNMI files have description text before the actual data
		if len(line) == 0 || (line[0] < '0' || line[0] > '9') {
			continue
		}

		record, err := parseLine(line, lineNum)
		if err != nil {
			return nil, err
		}

		records = append(records, *record)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading input: %w", err)
	}

	return records, nil
}

// parseLine parses a single data line.
func parseLine(line string, lineNum int) (*WeatherRecord, error) {
	fields := strings.Split(line, ",")

	if len(fields) != ExpectedColumns {
		return nil, fmt.Errorf("line %d: expected %d columns, got %d", lineNum, ExpectedColumns, len(fields))
	}

	// Parse station ID (required)
	stationID, err := parseRequiredInt(fields[0], "station_id", lineNum)
	if err != nil {
		return nil, err
	}

	// Parse date (required)
	date, err := parseDate(fields[1], lineNum)
	if err != nil {
		return nil, err
	}

	record := &WeatherRecord{
		StationID: stationID,
		Date:      date,
	}

	// Parse optional integer fields
	record.DDVEC = parseOptionalInt(fields[2])
	record.FHVEC = parseOptionalInt(fields[3])
	record.FG = parseOptionalInt(fields[4])
	record.FHX = parseOptionalInt(fields[5])
	record.FHXH = parseOptionalInt(fields[6])
	record.FHN = parseOptionalInt(fields[7])
	record.FHNH = parseOptionalInt(fields[8])
	record.FXX = parseOptionalInt(fields[9])
	record.FXXH = parseOptionalInt(fields[10])
	record.TG = parseOptionalInt(fields[11])
	record.TN = parseOptionalInt(fields[12])
	record.TNH = parseOptionalInt(fields[13])
	record.TX = parseOptionalInt(fields[14])
	record.TXH = parseOptionalInt(fields[15])
	record.T10N = parseOptionalInt(fields[16])
	record.T10NH = parseOptionalInt(fields[17])
	record.SQ = parseOptionalInt(fields[18])
	record.SP = parseOptionalInt(fields[19])
	record.Q = parseOptionalInt(fields[20])
	record.DR = parseOptionalInt(fields[21])
	record.RH = parseOptionalInt(fields[22])
	record.RHX = parseOptionalInt(fields[23])
	record.RHXH = parseOptionalInt(fields[24])
	record.PG = parseOptionalInt(fields[25])
	record.PX = parseOptionalInt(fields[26])
	record.PXH = parseOptionalInt(fields[27])
	record.PN = parseOptionalInt(fields[28])
	record.PNH = parseOptionalInt(fields[29])
	record.VVN = parseOptionalInt(fields[30])
	record.VVNH = parseOptionalInt(fields[31])
	record.VVX = parseOptionalInt(fields[32])
	record.VVXH = parseOptionalInt(fields[33])
	record.NG = parseOptionalInt(fields[34])
	record.UG = parseOptionalInt(fields[35])
	record.UX = parseOptionalInt(fields[36])
	record.UXH = parseOptionalInt(fields[37])
	record.UN = parseOptionalInt(fields[38])
	record.UNH = parseOptionalInt(fields[39])
	record.EV24 = parseOptionalInt(fields[40])

	return record, nil
}

// parseRequiredInt parses a required integer field.
func parseRequiredInt(s, name string, lineNum int) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("line %d: required field %s is empty", lineNum, name)
	}

	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("line %d: invalid %s value %q: %w", lineNum, name, s, err)
	}

	return v, nil
}

// parseOptionalInt parses an optional integer field, returning nil if empty.
func parseOptionalInt(s string) *int {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}

	v, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}

	return &v
}

// parseDate parses a date in YYYYMMDD format.
func parseDate(s string, lineNum int) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, fmt.Errorf("line %d: date is empty", lineNum)
	}

	t, err := time.Parse("20060102", s)
	if err != nil {
		return time.Time{}, fmt.Errorf("line %d: invalid date format %q (expected YYYYMMDD): %w", lineNum, s, err)
	}

	return t, nil
}
