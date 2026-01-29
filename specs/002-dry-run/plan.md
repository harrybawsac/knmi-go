# Implementation Plan: Dry-Run Mode for Sync Command

**Branch**: `002-dry-run` | **Date**: 2026-01-28 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/002-dry-run/spec.md`

## Summary

Add a `--dry-run` flag to the `knmi sync` command that downloads and parses KNMI weather data, determines which records would be inserted (excluding duplicates), and displays the last 10 records in a tabular format without actually inserting them into the database.

## Technical Context

**Language/Version**: Go 1.25  
**Primary Dependencies**: github.com/spf13/cobra (CLI), github.com/lib/pq (PostgreSQL)  
**Storage**: PostgreSQL (existing weather_records table)  
**Testing**: `go test` with standard library  
**Target Platform**: macOS/Linux CLI  
**Project Type**: Single CLI application  
**Performance Goals**: Preview results in <30 seconds (download + parse + display)  
**Constraints**: Must not insert any records when --dry-run is enabled  
**Scale/Scope**: Small feature addition to existing sync command

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Code Quality | ✅ PASS | Will use go vet, staticcheck, gofmt |
| II. Testing Standards | ✅ PASS | Will add unit tests for new functions, table-driven tests |
| III. CLI Design | ✅ PASS | Follows Unix conventions, --dry-run is standard flag |

**Gate Status**: ✅ PASSED - No violations

## Project Structure

### Documentation (this feature)

```text
specs/002-dry-run/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output (minimal - no new entities)
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
│   └── cli.md           # Updated CLI contract for --dry-run flag
└── tasks.md             # Phase 2 output
```

### Source Code (repository root)

```text
internal/
├── cli/
│   └── sync.go          # MODIFY: Add --dry-run flag and preview logic
├── db/
│   └── weather.go       # MODIFY: Add FilterNewRecords method
└── parser/
    └── csv.go           # EXISTING: No changes needed
```

**Structure Decision**: Minimal changes to existing structure. Modifications to sync.go and weather.go only.

## Complexity Tracking

> No violations - table left empty.

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| - | - | - |
