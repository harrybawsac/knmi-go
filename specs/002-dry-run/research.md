# Research: Dry-Run Mode for Sync Command

**Feature**: 002-dry-run | **Date**: 2026-01-28

## Overview

This feature has minimal research requirements as it builds entirely on existing infrastructure. All decisions leverage established patterns in the codebase.

## Research Tasks

### R-001: Duplicate Detection Strategy

**Question**: How to identify which records would be new (not already in DB) without inserting?

**Decision**: Query existing records by composite key (station_id, date) and filter in Go

**Rationale**: 
- The existing `InsertRecords()` uses `ON CONFLICT DO NOTHING` for duplicate handling
- For dry-run, we need to know *which* records are new before "inserting"
- Query for existing (station_id, date) pairs, filter parsed records against this set
- This is the same approach the DB would use, just explicit in application code

**Alternatives Considered**:
1. ~~Use PostgreSQL `INSERT ... RETURNING` with `ON CONFLICT DO NOTHING`~~ - Would still insert records
2. ~~Use temporary table~~ - Unnecessary complexity for preview feature
3. ~~Just show last 10 parsed records without dedup~~ - Doesn't match spec (must show what *would* be inserted)

### R-002: Tabular Output Format

**Question**: How to format weather records as a readable table in terminal?

**Decision**: Use fixed-width columns with fmt.Printf formatting

**Rationale**:
- Go's fmt package handles column alignment natively
- No external dependencies needed
- Simple format string: `%-10s %-10s %6s %6s %6s %6s %6s\n`
- Matches constitution principle of minimal dependencies

**Alternatives Considered**:
1. ~~Use tablewriter library~~ - External dependency for simple use case
2. ~~CSV output~~ - Less readable in terminal
3. ~~JSON output~~ - Not requested in spec

### R-003: Cobra Flag Pattern

**Question**: How to add --dry-run flag to existing sync command?

**Decision**: Add bool flag with `cmd.Flags().BoolVarP()`

**Rationale**:
- Matches existing flag patterns in sync.go (uses Cobra)
- Standard Cobra convention for boolean flags
- Short flag `-n` follows common convention (e.g., make -n, rsync -n)

**Implementation**:
```go
var dryRun bool
syncCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "Preview records without inserting")
```

## Summary

All research tasks resolved. No external research required. All decisions leverage existing patterns in the codebase.

| Task | Decision | Confidence |
|------|----------|------------|
| R-001 | Query + filter approach | High |
| R-002 | fmt.Printf fixed-width | High |
| R-003 | BoolVarP with -n short | High |
