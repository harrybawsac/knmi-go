# Tasks: Dry-Run Mode for Sync Command

**Input**: Design documents from `/specs/002-dry-run/`
**Prerequisites**: plan.md ✓, spec.md ✓, research.md ✓, data-model.md ✓, contracts/cli.md ✓

**Tests**: Not explicitly requested - omitting test tasks per task generation rules.

**Organization**: Tasks grouped by user story for independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1)
- Include exact file paths in descriptions

---

## Phase 1: Setup

**Purpose**: No setup tasks required - this feature extends an existing CLI project with established structure.

*(Skip to Phase 2)*

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Database method that must exist before sync command can filter new records

**⚠️ CRITICAL**: User Story 1 cannot proceed without the FilterNewRecords method

- [X] T001 Add FilterNewRecords method to WeatherRepository in internal/db/weather.go

**Checkpoint**: Foundation ready - FilterNewRecords method available for dry-run preview logic

---

## Phase 3: User Story 1 - Preview Data Before Sync (Priority: P1) 🎯 MVP

**Goal**: Allow users to preview weather records that would be inserted without modifying the database

**Independent Test**: Run `knmi sync --dry-run` and verify last 10 records are displayed without database modifications

### Implementation for User Story 1

- [X] T002 [US1] Add --dry-run flag to sync command in internal/cli/sync.go
- [X] T003 [US1] Add formatPreviewValue helper function for null handling in internal/cli/sync.go
- [X] T004 [US1] Add printPreviewTable function for tabular output (including empty result message per FR-008) in internal/cli/sync.go
- [X] T005 [US1] Add dry-run branch logic in runSync function in internal/cli/sync.go

**Checkpoint**: User Story 1 complete - `knmi sync --dry-run` shows preview without inserting

---

## Phase 4: Polish & Cross-Cutting Concerns

**Purpose**: Validation and documentation

- [X] T006 Run go vet and staticcheck on modified files
- [X] T007 Run gofmt on modified files
- [X] T008 Validate dry-run output matches quickstart.md examples
- [X] T009 Verify normal sync mode unchanged (regression check)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 2 (Foundational)**: No dependencies - can start immediately
- **Phase 3 (User Story 1)**: Depends on T001 (FilterNewRecords method)
- **Phase 4 (Polish)**: Depends on Phase 3 completion

### Task Dependencies within User Story 1

```
T001 (FilterNewRecords) ──┐
                          ├──► T002 (--dry-run flag)
                          │        │
                          │        ▼
                          ├──► T003 (formatPreviewValue helper)
                          │        │
                          │        ▼
                          └──► T004 (printPreviewTable function)
                                   │
                                   ▼
                              T005 (dry-run branch in runSync)
```

### Parallel Opportunities

- T003 and T004 can be developed in parallel (helper functions, no dependencies between them)
- All T006-T009 polish tasks can run in parallel

---

## Parallel Example: User Story 1

```bash
# After T002 completes, launch helper functions together:
Task: "T003 [US1] Add formatPreviewValue helper function"
Task: "T004 [US1] Add printPreviewTable function"

# After helpers complete:
Task: "T005 [US1] Add dry-run branch logic in runSync"
```

---

## Implementation Strategy

### MVP First (This is a single-story feature)

1. Complete Phase 2: Foundational (T001)
2. Complete Phase 3: User Story 1 (T002-T005)
3. **STOP and VALIDATE**: Test with `knmi sync --dry-run`
4. Complete Phase 4: Polish (T006-T009)

### Incremental Delivery

1. T001 → FilterNewRecords method ready
2. T002 → --dry-run flag recognized
3. T003-T004 → Helper functions for output
4. T005 → Full dry-run flow working
5. T006-T009 → Code quality validation

---

## Notes

- Single user story feature: P1 "Preview Data Before Sync" covers all requirements
- No schema changes required per data-model.md
- FilterNewRecords queries existing (station_id, date) pairs and filters in Go per research.md
- Output format uses fmt.Printf fixed-width columns per research.md
- All tasks modify only 2 files: internal/cli/sync.go and internal/db/weather.go
