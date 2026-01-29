# Data Model: Dry-Run Mode for Sync Command

**Feature**: 002-dry-run | **Date**: 2026-01-28

## Overview

This feature requires **no schema changes**. The dry-run mode operates on the existing `weather_records` table and `parser.WeatherRecord` struct.

## Existing Entities (Reference)

### weather_records (PostgreSQL Table)

```sql
CREATE TABLE weather_records (
    id SERIAL PRIMARY KEY,
    station_id INTEGER NOT NULL,
    date DATE NOT NULL,
    tg INTEGER,      -- Mean temperature (0.1°C)
    tn INTEGER,      -- Minimum temperature (0.1°C)
    tx INTEGER,      -- Maximum temperature (0.1°C)
    fg INTEGER,      -- Mean wind speed (0.1 m/s)
    rh INTEGER,      -- Precipitation (0.1 mm)
    -- ... 34 additional weather measurement fields
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(station_id, date)
);
```

### parser.WeatherRecord (Go Struct)

```go
type WeatherRecord struct {
    StationID int        // STN field
    Date      time.Time  // YYYYMMDD field
    TG        *int       // Mean temperature
    TN        *int       // Minimum temperature
    TX        *int       // Maximum temperature
    FG        *int       // Mean wind speed
    RH        *int       // Precipitation
    // ... 34 additional fields
}
```

## New Types

### PreviewRecord (Display Struct)

A subset of WeatherRecord for dry-run output (7 fields per FR-006):

```go
// PreviewRecord contains fields displayed in dry-run preview
type PreviewRecord struct {
    Date      string  // Formatted as YYYY-MM-DD
    StationID int     // Station identifier
    TG        string  // Mean temp (formatted with unit handling)
    TN        string  // Min temp
    TX        string  // Max temp
    FG        string  // Wind speed
    RH        string  // Precipitation
}
```

**Note**: This is an internal display struct, not a database entity. Fields are strings for formatted output with null handling (displays "-" for nil values).

## Relationships

No new relationships. Dry-run mode reads from `weather_records` to determine existing records.

## Validation Rules

No new validation rules. Existing parsing validation in `internal/parser/csv.go` applies.

## State Transitions

N/A - This feature is stateless (read-only operation).
