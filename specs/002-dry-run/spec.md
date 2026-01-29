# Feature Specification: Dry-Run Mode for Sync Command

**Feature Branch**: `002-dry-run`  
**Created**: 2026-01-28  
**Status**: Draft  
**Input**: User description: "Add --dry-run flag to sync command that fetches KNMI data and outputs the last 10 lines it would INSERT to the console without actually inserting"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Preview Data Before Sync (Priority: P1)

As a user, I want to preview the weather records that would be inserted into the database before actually inserting them, so I can verify the data looks correct and understand what changes will be made.

**Why this priority**: This is the core value of the feature - allowing users to preview data before committing to database changes. It enables safe data validation and builds user confidence.

**Independent Test**: Run `knmi sync --dry-run` and verify that the last 10 records to be inserted are displayed on the console without any database modifications.

**Acceptance Scenarios**:

1. **Given** a configured database URL and KNMI data source, **When** the user runs `knmi sync --dry-run`, **Then** the system downloads and parses KNMI data and displays the last 10 records that would be inserted without modifying the database.

2. **Given** a database with existing records, **When** the user runs `knmi sync --dry-run`, **Then** only new records (not duplicates) are shown in the preview output.

3. **Given** a successful dry-run execution, **When** checking the database, **Then** no new records have been inserted.

---

### Edge Cases

- What happens when there are fewer than 10 new records? Display all available new records.
- What happens when there are no new records to insert? Display a message indicating no new data available.
- What happens when database connection fails? Display an error (database check is needed to determine duplicates).
- What happens when the KNMI data source is unavailable? Display a download error.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST accept a `--dry-run` flag on the `sync` command
- **FR-002**: System MUST download and parse KNMI data in dry-run mode (same as normal sync)
- **FR-003**: System MUST connect to the database to determine which records are new (not duplicates)
- **FR-004**: System MUST display the last 10 records that would be inserted to stdout
- **FR-005**: System MUST NOT insert any records into the database when `--dry-run` is enabled
- **FR-006**: System MUST display record data in a tabular format with 7 fields: date, station_id, tg (mean temp), tn (min temp), tx (max temp), fg (wind speed), rh (precipitation)
- **FR-007**: System MUST indicate the total count of records that would be inserted
- **FR-008**: System MUST display a message when there are no new records to insert

### Key Entities

- **WeatherRecord**: Existing entity - represents a single daily weather observation with 41 data columns
- **Preview Output**: Human-readable representation of weather records showing key fields (date, station, temperature, etc.)

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can preview sync results in under 30 seconds (download + parse + display)
- **SC-002**: Running `knmi sync --dry-run` does not modify the database (zero rows inserted)
- **SC-003**: Users can see exactly which records would be inserted before committing
- **SC-004**: Preview output shows 7 fields: date, station_id, tg, tn, tx, fg, rh in tabular format

## Assumptions

- The existing sync infrastructure (download, parse, duplicate detection) will be reused
- Database connection is required to check for existing records and determine true new records
- The 10-record limit is sufficient for preview purposes (avoids overwhelming output)
- Output format should be tabular or structured for readability

## Clarifications

### Session 2026-01-28

- Q: Which fields to display in preview output? → A: Standard 7 fields: date, station_id, tg (temp), tn (min temp), tx (max temp), fg (wind), rh (rain)
