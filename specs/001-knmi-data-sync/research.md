# Research: KNMI Weather Data Sync CLI

**Date**: 2026-01-28  
**Feature**: 001-knmi-data-sync

## KNMI Data Format

### File Structure

- **Source URL**: `https://cdn.knmi.nl/knmi/map/page/klimatologie/gegevens/daggegevens/etmgeg_260.zip`
- **Archive**: ZIP file containing a single `.txt` file (`etmgeg_260.txt`)
- **Format**: CSV-like with semicolon-delimited comments, comma-delimited data
- **Header lines**: ~50 lines of metadata/documentation before data starts
- **Header row**: Starts with `# STN,YYYYMMDD,...` (line with column names)
- **Data rows**: Start with station number (e.g., `  260,19010101,...`)
- **Total rows**: ~45,700 records (data from 1901-01-01 to present)
- **Values**: Integers, often in 0.1 units (e.g., temperature in 0.1°C), empty cells for missing data

### Column Definitions (41 columns)

| Column | Description | Unit |
|--------|-------------|------|
| STN | Station number | integer |
| YYYYMMDD | Date | YYYYMMDD format |
| DDVEC | Vector mean wind direction | degrees (360=N, 90=E, 180=S, 270=W, 0=calm) |
| FHVEC | Vector mean windspeed | 0.1 m/s |
| FG | Daily mean windspeed | 0.1 m/s |
| FHX | Maximum hourly mean windspeed | 0.1 m/s |
| FHXH | Hour of FHX | hour (1-24) |
| FHN | Minimum hourly mean windspeed | 0.1 m/s |
| FHNH | Hour of FHN | hour (1-24) |
| FXX | Maximum wind gust | 0.1 m/s |
| FXXH | Hour of FXX | hour (1-24) |
| TG | Daily mean temperature | 0.1 °C |
| TN | Minimum temperature | 0.1 °C |
| TNH | Hour of TN | hour (1-24) |
| TX | Maximum temperature | 0.1 °C |
| TXH | Hour of TX | hour (1-24) |
| T10N | Minimum temp at 10cm | 0.1 °C |
| T10NH | 6-hour period of T10N | 6=0-6, 12=6-12, 18=12-18, 24=18-24 UT |
| SQ | Sunshine duration | 0.1 hour (-1 for <0.05h) |
| SP | Sunshine percentage | percent |
| Q | Global radiation | J/cm² |
| DR | Precipitation duration | 0.1 hour |
| RH | Daily precipitation | 0.1 mm (-1 for <0.05mm) |
| RHX | Maximum hourly precipitation | 0.1 mm (-1 for <0.05mm) |
| RHXH | Hour of RHX | hour (1-24) |
| PG | Daily mean sea level pressure | 0.1 hPa |
| PX | Maximum hourly pressure | 0.1 hPa |
| PXH | Hour of PX | hour (1-24) |
| PN | Minimum hourly pressure | 0.1 hPa |
| PNH | Hour of PN | hour (1-24) |
| VVN | Minimum visibility | coded (0=<100m, 50=5-6km, 89=>70km) |
| VVNH | Hour of VVN | hour (1-24) |
| VVX | Maximum visibility | coded |
| VVXH | Hour of VVX | hour (1-24) |
| NG | Mean cloud cover | octants (0-8, 9=invisible) |
| UG | Daily mean relative humidity | percent |
| UX | Maximum relative humidity | percent |
| UXH | Hour of UX | hour (1-24) |
| UN | Minimum relative humidity | percent |
| UNH | Hour of UN | hour (1-24) |
| EV24 | Potential evapotranspiration | 0.1 mm |

### Parsing Considerations

1. **Skip metadata**: Lines not starting with station number should be skipped
2. **Trim whitespace**: Values have leading/trailing spaces
3. **Empty values**: Missing data represented as empty strings (` , ,`)
4. **Integer conversion**: All values are integers; empty → NULL in database
5. **Date parsing**: YYYYMMDD format → Go `time.Time`

## Technology Decisions

### PostgreSQL Driver

**Decision**: Use `github.com/lib/pq`  
**Rationale**: Most widely used, pure Go PostgreSQL driver. Stable, well-documented.  
**Alternatives considered**:
- `pgx`: More performant but more complex. Overkill for this use case.
- `database/sql` only: Would need a driver anyway.

### CLI Framework

**Decision**: Use `github.com/spf13/cobra`  
**Rationale**: De facto standard for Go CLIs. Provides subcommands, flags, help generation.  
**Alternatives considered**:
- Standard library `flag`: No subcommand support, would need manual help formatting.
- `urfave/cli`: Good but less widely adopted than Cobra.

### Migration Library

**Decision**: Custom minimal implementation  
**Rationale**: Feature spec requires reading SQL files from a directory and tracking applied migrations. A simple tracker table + file reader is sufficient (~100 LOC). External libraries add complexity.  
**Alternatives considered**:
- `golang-migrate/migrate`: Full-featured but heavyweight for simple sequential migrations.
- `pressly/goose`: Similar, adds unnecessary dependency for this scope.

### Configuration

**Decision**: Environment variables with flag overrides  
**Rationale**: Standard 12-factor app approach. `DATABASE_URL` for connection string.  
**Alternatives considered**:
- Config files: Unnecessary complexity for a simple CLI.
- Viper: Overkill for 2-3 config values.

## Best Practices Applied

### Go CLI Best Practices

1. **Exit codes**: 0 for success, 1 for errors (aligns with Constitution III)
2. **Output streams**: stdout for results, stderr for errors and verbose logs
3. **Flags**: Use long names (`--verbose` not `-v` alone) for clarity
4. **Help**: Auto-generated via Cobra, includes examples

### Database Best Practices

1. **Transactions**: All migrations run in transactions for atomicity
2. **Idempotent sync**: Use `INSERT ... ON CONFLICT DO NOTHING` for upserts
3. **Connection pooling**: Let `database/sql` handle connection pool
4. **Indexes**: Create index on (station_id, date) for unique constraint and fast lookups

### Testing Strategy

1. **Unit tests**: Parser, migration file discovery, config parsing
2. **Integration tests**: Use Docker PostgreSQL for real database tests
3. **Test fixtures**: Include sample KNMI data in testdata/ directory
4. **Table-driven tests**: For parser with various edge cases (empty values, negative numbers)

## Resolved Clarifications

| Question | Resolution |
|----------|------------|
| Store all KNMI columns? | Yes, store all 41 columns (comprehensive) |
| Logging behavior? | Quiet by default, `--verbose` flag for progress |
| Unique key for records? | station_id + date combination |
| Handle missing values? | Store as NULL in database |
