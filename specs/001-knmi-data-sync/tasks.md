# Tasks: KNMI Weather Data Sync CLI

**Input**: Design documents from `/specs/001-knmi-data-sync/`
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, contracts/ ✅

**Tests**: Included per Constitution II (Testing Standards) — unit tests for exported functions, table-driven tests for parser.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Single project**: `cmd/`, `internal/`, `migrations/`, `tests/` at repository root

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Initialize Go module with `go mod init` in repository root
- [x] T002 [P] Create project directory structure per plan.md (cmd/, internal/, migrations/, tests/)
- [x] T003 [P] Add .golangci.yml for linting configuration
- [x] T004 [P] Add Makefile with build, test, lint, fmt targets

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T005 Implement configuration loader in internal/config/config.go (DATABASE_URL, KNMI_DATA_URL, KNMI_MIGRATIONS_DIR)
- [x] T006 [P] Implement database connection pool in internal/db/connection.go
- [x] T007 [P] Create initial SQL migration file in migrations/001_create_tables.sql (weather_records + migrations tables per data-model.md)
- [x] T008 Add cobra dependency and create root command skeleton in internal/cli/root.go (--help, --version, --verbose, --database-url flags)
- [x] T009 Create CLI entrypoint in cmd/knmi/main.go

**Checkpoint**: Foundation ready — `knmi --help` and `knmi --version` work

---

## Phase 3: User Story 1 - Run Database Migrations (Priority: P1) 🎯 MVP

**Goal**: Users can run `knmi migrate` to apply SQL migrations and set up the database schema.

**Independent Test**: Run `knmi migrate` with a fresh database; verify tables are created.

### Tests for User Story 1

- [x] T010 [P] [US1] Unit tests for migration file discovery and version parsing (table-driven) in tests/unit/migration_test.go
- [x] T011 [P] [US1] Integration test for migrate command in tests/integration/migrate_test.go

### Implementation for User Story 1

- [x] T012 [US1] Implement migration file discovery (read/sort by version) in internal/migration/runner.go
- [x] T013 [US1] Implement migration state tracker (query/insert applied migrations) in internal/migration/tracker.go
- [x] T014 [US1] Implement migration execution with transactions in internal/migration/runner.go
- [x] T015 [US1] Create migrate command in internal/cli/migrate.go (--migrations-dir flag)
- [x] T016 [US1] Add error handling for invalid SQL, connection failures per contracts/cli.md
- [x] T017 [US1] Add verbose logging for migration progress

**Checkpoint**: `knmi migrate` works — schema is created, migrations are tracked

---

## Phase 4: User Story 2 - Sync Weather Data (Priority: P2)

**Goal**: Users can run `knmi sync` to download KNMI data and insert new weather records.

**Independent Test**: Run `knmi sync` after migrations; verify records appear in database.

### Tests for User Story 2

- [x] T018 [P] [US2] Unit test for CSV parsing (table-driven, edge cases) in tests/unit/parser_test.go
- [x] T019 [P] [US2] Unit tests for HTTP downloader and zip extraction in tests/unit/fetch_test.go
- [x] T020 [P] [US2] Integration test for sync command in tests/integration/sync_test.go

### Implementation for User Story 2

- [x] T021 [P] [US2] Implement HTTP downloader in internal/fetch/downloader.go
- [x] T022 [P] [US2] Implement zip extractor in internal/fetch/extractor.go
- [x] T023 [US2] Implement KNMI CSV parser (41 columns, handle empty values, warn on unexpected columns) in internal/parser/csv.go
- [x] T024 [US2] Implement weather record repository (insert with ON CONFLICT DO NOTHING) in internal/db/weather.go
- [x] T025 [US2] Create sync command in internal/cli/sync.go (--url flag)
- [x] T026 [US2] Add "no migrations applied" check before sync
- [x] T027 [US2] Add error handling for network, parse, database errors per contracts/cli.md
- [x] T028 [US2] Add verbose logging for download/parse/insert progress

**Checkpoint**: `knmi sync` works — weather data is downloaded, parsed, and inserted

---

## Phase 5: User Story 3 - Display Version and Help (Priority: P3)

**Goal**: Users can see CLI help and version information.

**Independent Test**: Run `knmi --help` and `knmi --version`; verify correct output.

### Implementation for User Story 3

- [x] T029 [US3] Add version command in internal/cli/root.go
- [x] T030 [US3] Add help text with examples for all commands per contracts/cli.md
- [x] T031 [US3] Add unknown command error handling with suggestions

**Checkpoint**: CLI usability complete — help, version, and error guidance work

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final cleanup and validation

- [x] T032 [P] Run `go vet ./...` and fix any warnings
- [x] T033 [P] Run `staticcheck ./...` and fix any warnings
- [x] T034 [P] Run `gofmt -w .` to format all code
- [x] T035 [P] Add README.md with installation and usage instructions
- [x] T036 Verify test coverage ≥70% with `go test -cover ./...`
- [x] T037 Run quickstart.md validation (end-to-end test)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion — BLOCKS all user stories
- **User Stories (Phase 3-5)**: All depend on Foundational phase completion
  - User stories can proceed sequentially in priority order (P1 → P2 → P3)
  - US2 depends on US1 (needs migrations to exist)
  - US3 is independent but best done after US1/US2 for complete help text
- **Polish (Phase 6)**: Depends on all user stories being complete

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Infrastructure (downloader, parser) before commands
- Core implementation before error handling
- Error handling before verbose logging

### Parallel Opportunities

**Phase 1** (all can run in parallel):
```
T002, T003, T004
```

**Phase 2** (T006, T007 can run in parallel):
```
T006, T007
```

**User Story 1 Tests** (all can run in parallel):
```
T010, T011
```

**User Story 2 Tests** (all can run in parallel):
```
T018, T019, T020
```

**User Story 2 Implementation** (T021, T022 can run in parallel):
```
T021, T022
```

**Phase 6** (T032-T035 can run in parallel):
```
T032, T033, T034, T035
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: User Story 1 (Migrations)
4. **STOP and VALIDATE**: Test `knmi migrate` independently
5. Deploy/demo if schema management alone is useful

### Incremental Delivery

1. Complete Setup + Foundational → `knmi --version` works
2. Add User Story 1 → `knmi migrate` works → Validate
3. Add User Story 2 → `knmi sync` works → Validate
4. Add User Story 3 → Help/version polished → Validate
5. Polish → Production-ready

---

## Summary

| Metric | Value |
|--------|-------|
| **Total Tasks** | 37 |
| **Phase 1 (Setup)** | 4 tasks |
| **Phase 2 (Foundational)** | 5 tasks |
| **User Story 1 (Migrations)** | 8 tasks (2 tests + 6 impl) |
| **User Story 2 (Sync)** | 11 tasks (3 tests + 8 impl) |
| **User Story 3 (Help/Version)** | 3 tasks |
| **Phase 6 (Polish)** | 6 tasks |
| **Parallel opportunities** | 17 tasks marked [P] |
| **MVP scope** | Phase 1 + 2 + 3 (17 tasks) |
