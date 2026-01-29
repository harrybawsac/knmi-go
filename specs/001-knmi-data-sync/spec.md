# Feature Specification: KNMI Weather Data Sync CLI

**Feature Branch**: `001-knmi-data-sync`  
**Created**: 2026-01-28  
**Status**: Draft  
**Input**: User description: "Build a CLI that fetches KNMI weather data from a zip file, parses the CSV, and syncs it to PostgreSQL with migration support"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Run Database Migrations (Priority: P1)

As a user, I want to run SQL migration files against my PostgreSQL database so that the schema is properly set up before importing weather data.

**Why this priority**: Without a proper database schema, no data can be stored. Migrations are the foundation that all other features depend on.

**Independent Test**: Can be fully tested by running `knmi migrate` with a fresh database and verifying tables are created. Delivers value by ensuring repeatable, version-controlled schema management.

**Acceptance Scenarios**:

1. **Given** a directory with numbered SQL migration files (e.g., `001_create_tables.sql`), **When** I run `knmi migrate`, **Then** all pending migrations are applied in order and the migration state is tracked.
2. **Given** some migrations have already been applied, **When** I run `knmi migrate`, **Then** only new migrations are applied.
3. **Given** a migration file contains invalid SQL, **When** I run `knmi migrate`, **Then** the migration fails with a clear error message and no partial changes are committed.

---

### User Story 2 - Sync Weather Data (Priority: P2)

As a user, I want to fetch KNMI weather data and sync it to my database so that I have up-to-date historical weather records.

**Why this priority**: This is the core functionality of the CLI—fetching, parsing, and storing weather data. Depends on P1 (migrations) being complete.

**Independent Test**: Can be fully tested by running `knmi sync` after migrations and verifying new records appear in the database. Delivers value by automating weather data ingestion.

**Acceptance Scenarios**:

1. **Given** a database with the weather schema in place, **When** I run `knmi sync`, **Then** the CLI downloads the KNMI zip file, extracts the CSV data, and inserts new records.
2. **Given** the database already contains some weather records, **When** I run `knmi sync`, **Then** only missing records are inserted (no duplicates).
3. **Given** the KNMI server is unreachable, **When** I run `knmi sync`, **Then** the CLI exits with an error message indicating the network issue.
4. **Given** the zip file or CSV format is corrupted, **When** I run `knmi sync`, **Then** the CLI exits with a descriptive parse error.

---

### User Story 3 - Display Version and Help (Priority: P3)

As a user, I want to see help information and the CLI version so that I understand how to use the tool and which version I'm running.

**Why this priority**: Standard CLI usability features. Lower priority because the tool is functional without them, but essential for a polished user experience.

**Independent Test**: Can be fully tested by running `knmi --help` and `knmi --version` and verifying correct output. Delivers value by providing discoverability and debugging context.

**Acceptance Scenarios**:

1. **Given** the CLI is installed, **When** I run `knmi --help`, **Then** I see a list of available commands and flags with descriptions.
2. **Given** the CLI is installed, **When** I run `knmi --version`, **Then** I see the version number (e.g., `knmi v1.0.0`).
3. **Given** I run an unknown command, **When** the CLI processes it, **Then** I see an error message with suggestions and a reference to `--help`.

---

### Edge Cases

- What happens when the database connection fails mid-migration? → Transaction is rolled back, no partial state.
- What happens when the zip file URL changes? → URL should be configurable via flag or environment variable.
- What happens when the CSV has unexpected columns or format changes? → CLI logs a warning but attempts to parse known columns.
- What happens when running sync with an empty database (no migrations)? → CLI exits with an error prompting user to run migrations first.
- What happens with concurrent sync operations? → Database constraints prevent duplicate records; CLI handles gracefully.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: CLI MUST provide a `migrate` command that applies SQL migration files in sequential order.
- **FR-002**: CLI MUST track which migrations have been applied to avoid re-running them.
- **FR-003**: CLI MUST provide a `sync` command that downloads the KNMI zip file from a configurable URL.
- **FR-004**: CLI MUST extract and parse the `.txt` file from the zip as CSV data.
- **FR-005**: CLI MUST insert only new weather records using `(station_id, date)` as unique key; duplicates are skipped via `INSERT ... ON CONFLICT DO NOTHING`.
- **FR-006**: CLI MUST support `--help` and `--version` flags.
- **FR-011**: CLI MUST be quiet by default (errors only); a `--verbose` flag enables detailed progress logging.
- **FR-007**: CLI MUST write errors to stderr and normal output to stdout.
- **FR-008**: CLI MUST exit with code 0 on success and non-zero on failure.
- **FR-009**: CLI MUST accept database connection details via environment variable or flag.
- **FR-010**: CLI MUST use transactions for migrations to ensure atomicity.

### Key Entities

- **WeatherRecord**: A single day's weather observation from KNMI. Stores all available columns from the KNMI CSV (comprehensive approach). Common attributes include: station number, date, temperature metrics, precipitation, wind measurements, sunshine duration, humidity, pressure, and other meteorological data provided by KNMI.
- **Migration**: A versioned SQL file. Key attributes: sequence number, filename, applied timestamp.
- **Station**: Weather station metadata. Key attributes: station number, name, location (implicit from KNMI data).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can set up the database schema by running a single command (`knmi migrate`) in under 10 seconds.
- **SC-002**: Users can sync the latest KNMI data by running a single command (`knmi sync`) in under 60 seconds for initial load.
- **SC-003**: Subsequent sync operations complete in under 10 seconds when no new data exists.
- **SC-004**: CLI provides clear, actionable error messages for all failure scenarios (network, parse, database errors).
- **SC-005**: All commands are documented via `--help` with examples.

## Clarifications

### Session 2026-01-28

- Q: Should the CLI store all available KNMI columns or only the 6 listed attributes? → A: Store all available KNMI columns (comprehensive).
- Q: Should the CLI log progress during normal operations? → A: Quiet by default, `--verbose` flag for detailed progress.

## Assumptions

- The KNMI zip file format (containing a `.txt` file with CSV-like structure) will remain stable.
- PostgreSQL is the target database; no other databases need to be supported initially.
- The default data source URL is `https://cdn.knmi.nl/knmi/map/page/klimatologie/gegevens/daggegevens/etmgeg_260.zip`.
- Migration files follow the naming convention `NNN_description.sql` (e.g., `001_create_tables.sql`).
- Weather records are uniquely identified by station number + date combination.
