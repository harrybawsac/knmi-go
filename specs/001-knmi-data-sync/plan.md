# Implementation Plan: KNMI Weather Data Sync CLI

**Branch**: `001-knmi-data-sync` | **Date**: 2026-01-28 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/001-knmi-data-sync/spec.md`

## Summary

Build a Go CLI tool (`knmi`) that fetches weather data from the KNMI (Dutch meteorological institute) API, parses CSV data from a zip archive, and syncs it to a PostgreSQL database. The CLI includes a migration system for schema management and supports incremental data sync to avoid duplicates.

## Technical Context

**Language/Version**: Go 1.25  
**Primary Dependencies**: Standard library (net/http, archive/zip, encoding/csv), lib/pq (PostgreSQL driver)  
**Storage**: PostgreSQL  
**Testing**: `go test` with standard library `testing` package  
**Target Platform**: macOS, Linux (CLI binary)  
**Project Type**: Single project  
**Performance Goals**: Initial sync <60s, subsequent syncs <10s  
**Constraints**: Minimal external dependencies, quiet by default  
**Scale/Scope**: Single weather station data (station 260), ~50 years of daily records (~18,000 rows)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Pre-Research Check (2026-01-28)

| Principle | Requirement | Status |
|-----------|-------------|--------|
| I. Code Quality | Pass `go vet`, `staticcheck`, `gofmt`; handle all errors | ✅ PLANNED |
| I. Code Quality | Meaningful names, no ignored errors | ✅ PLANNED |
| II. Testing Standards | Unit tests for exported functions | ✅ PLANNED |
| II. Testing Standards | Table-driven tests for multi-scenario functions | ✅ PLANNED |
| II. Testing Standards | Tests runnable via `go test ./...` | ✅ PLANNED |
| II. Testing Standards | ≥70% coverage for new packages | ✅ PLANNED |
| III. CLI Design | Exit 0 success, non-zero errors | ✅ Matches FR-008 |
| III. CLI Design | Errors to stderr, output to stdout | ✅ Matches FR-007 |
| III. CLI Design | Support `--help` and `--version` | ✅ Matches FR-006 |
| III. CLI Design | Actionable error messages | ✅ Matches SC-004 |

**Pre-Research Gate Status**: ✅ PASS — No violations detected.

### Post-Design Check (2026-01-28)

| Principle | Design Element | Status |
|-----------|----------------|--------|
| I. Code Quality | Using `lib/pq` (stable, well-maintained) | ✅ PASS |
| I. Code Quality | Using `cobra` for CLI (idiomatic Go) | ✅ PASS |
| I. Code Quality | All errors explicitly handled in contracts | ✅ PASS |
| II. Testing Standards | Test structure defined in plan | ✅ PASS |
| II. Testing Standards | Table-driven tests planned for parser | ✅ PASS |
| III. CLI Design | Exit codes documented in contracts | ✅ PASS |
| III. CLI Design | stderr/stdout separation in contracts | ✅ PASS |
| III. CLI Design | `--help`, `--version` in contracts | ✅ PASS |
| III. CLI Design | Error message format specified | ✅ PASS |

**Post-Design Gate Status**: ✅ PASS — All design elements align with constitution.

## Project Structure

### Documentation (this feature)

```text
specs/001-knmi-data-sync/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (CLI interface spec)
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
cmd/
└── knmi/
    └── main.go          # CLI entrypoint

internal/
├── cli/
│   ├── root.go          # Root command, --help, --version
│   ├── migrate.go       # migrate command
│   └── sync.go          # sync command
├── migration/
│   ├── runner.go        # Migration execution logic
│   └── tracker.go       # Migration state tracking
├── fetch/
│   ├── downloader.go    # HTTP zip download
│   └── extractor.go     # Zip extraction
├── parser/
│   └── csv.go           # KNMI CSV parsing
├── db/
│   ├── connection.go    # PostgreSQL connection
│   └── weather.go       # Weather record operations
└── config/
    └── config.go        # Configuration (env vars, flags)

migrations/
└── 001_create_tables.sql

tests/
├── integration/
│   ├── migrate_test.go
│   └── sync_test.go
└── unit/
    ├── parser_test.go
    ├── migration_test.go
    └── fetch_test.go
```

**Structure Decision**: Single project with `cmd/` for entrypoint and `internal/` for packages (Go convention). Tests split by type. Migration SQL files in `migrations/` directory.

## Complexity Tracking

> No violations detected. Section left empty per template instructions.
