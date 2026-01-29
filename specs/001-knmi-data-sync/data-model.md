# Data Model: KNMI Weather Data Sync CLI

**Date**: 2026-01-28  
**Feature**: 001-knmi-data-sync

## Entities

### WeatherRecord

Stores daily weather observations from KNMI station 260 (De Bilt).

| Field | Type | Nullable | Description |
|-------|------|----------|-------------|
| id | SERIAL | No | Primary key |
| station_id | INTEGER | No | Station number (260) |
| date | DATE | No | Observation date |
| ddvec | INTEGER | Yes | Vector mean wind direction (degrees) |
| fhvec | INTEGER | Yes | Vector mean windspeed (0.1 m/s) |
| fg | INTEGER | Yes | Daily mean windspeed (0.1 m/s) |
| fhx | INTEGER | Yes | Max hourly windspeed (0.1 m/s) |
| fhxh | INTEGER | Yes | Hour of FHX |
| fhn | INTEGER | Yes | Min hourly windspeed (0.1 m/s) |
| fhnh | INTEGER | Yes | Hour of FHN |
| fxx | INTEGER | Yes | Max wind gust (0.1 m/s) |
| fxxh | INTEGER | Yes | Hour of FXX |
| tg | INTEGER | Yes | Daily mean temp (0.1 °C) |
| tn | INTEGER | Yes | Min temp (0.1 °C) |
| tnh | INTEGER | Yes | Hour of TN |
| tx | INTEGER | Yes | Max temp (0.1 °C) |
| txh | INTEGER | Yes | Hour of TX |
| t10n | INTEGER | Yes | Min temp at 10cm (0.1 °C) |
| t10nh | INTEGER | Yes | Period of T10N |
| sq | INTEGER | Yes | Sunshine duration (0.1 hour) |
| sp | INTEGER | Yes | Sunshine percentage |
| q | INTEGER | Yes | Global radiation (J/cm²) |
| dr | INTEGER | Yes | Precipitation duration (0.1 hour) |
| rh | INTEGER | Yes | Daily precipitation (0.1 mm) |
| rhx | INTEGER | Yes | Max hourly precipitation (0.1 mm) |
| rhxh | INTEGER | Yes | Hour of RHX |
| pg | INTEGER | Yes | Mean sea level pressure (0.1 hPa) |
| px | INTEGER | Yes | Max pressure (0.1 hPa) |
| pxh | INTEGER | Yes | Hour of PX |
| pn | INTEGER | Yes | Min pressure (0.1 hPa) |
| pnh | INTEGER | Yes | Hour of PN |
| vvn | INTEGER | Yes | Min visibility (coded) |
| vvnh | INTEGER | Yes | Hour of VVN |
| vvx | INTEGER | Yes | Max visibility (coded) |
| vvxh | INTEGER | Yes | Hour of VVX |
| ng | INTEGER | Yes | Mean cloud cover (octants) |
| ug | INTEGER | Yes | Mean relative humidity (%) |
| ux | INTEGER | Yes | Max relative humidity (%) |
| uxh | INTEGER | Yes | Hour of UX |
| un | INTEGER | Yes | Min relative humidity (%) |
| unh | INTEGER | Yes | Hour of UN |
| ev24 | INTEGER | Yes | Evapotranspiration (0.1 mm) |
| created_at | TIMESTAMP | No | Record creation timestamp |

**Constraints**:
- PRIMARY KEY: `id`
- UNIQUE: `(station_id, date)` — prevents duplicate records
- INDEX: On unique constraint for fast lookups

### Migration

Tracks applied database migrations.

| Field | Type | Nullable | Description |
|-------|------|----------|-------------|
| id | SERIAL | No | Primary key |
| version | INTEGER | No | Migration sequence number |
| name | VARCHAR(255) | No | Migration filename |
| applied_at | TIMESTAMP | No | When migration was applied |

**Constraints**:
- PRIMARY KEY: `id`
- UNIQUE: `version` — each migration applied once

## Relationships

```text
┌─────────────────┐
│    Migration    │  (standalone, tracks schema changes)
└─────────────────┘

┌─────────────────┐
│  WeatherRecord  │  (standalone, one per station+date)
└─────────────────┘
```

No foreign key relationships. Entities are independent.

## Validation Rules

### WeatherRecord

1. `station_id` must be a positive integer
2. `date` must be a valid date, not in the future
3. All measurement fields are nullable (missing data allowed)
4. Temperature values can be negative (valid for winter measurements)
5. Special value `-1` in SQ, RH, RHX means "less than 0.05 units" (store as-is)

### Migration

1. `version` must be unique and sequential
2. `name` must match the pattern `NNN_description.sql`
3. `applied_at` is set automatically on insertion

## State Transitions

### Migration States

```text
[Not Applied] → migrate command → [Applied]
```

Migrations are one-way. No rollback support in initial version.

### WeatherRecord States

```text
[Does Not Exist] → sync command → [Exists]
```

Records are immutable once inserted. Historical data does not change.
