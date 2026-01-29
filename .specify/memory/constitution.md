<!--
  SYNC IMPACT REPORT
  ==================
  Version change: N/A → 1.0.0 (initial ratification)
  
  Added Principles:
  - I. Code Quality
  - II. Testing Standards
  - III. CLI Design
  
  Added Sections:
  - Core Principles (3 principles)
  - Technical Stack
  - Governance
  
  Templates Status:
  - plan-template.md: ✅ Compatible (uses Constitution Check gate)
  - spec-template.md: ✅ Compatible (user stories align with testing principle)
  - tasks-template.md: ✅ Compatible (test tasks align with Testing Standards)
  
  Follow-up TODOs: None
-->

# knmi-go Constitution

## Core Principles

### I. Code Quality

All Go code MUST:
- Pass `go vet` and `staticcheck` with zero warnings
- Be formatted with `gofmt` or `goimports`
- Have no unused exports; unexported symbols preferred unless API requires export
- Use meaningful names: avoid single-letter variables except in small scopes (loops, lambdas)
- Handle all errors explicitly; no ignored error returns

**Rationale**: Go's simplicity is its strength. Clean, idiomatic code reduces cognitive load and catches bugs early.

### II. Testing Standards

All features MUST have corresponding tests:
- Unit tests for all exported functions and methods
- Table-driven tests preferred for functions with multiple input scenarios
- Tests MUST be runnable via `go test ./...` with no external dependencies
- Test coverage SHOULD be ≥70% for new packages; critical paths MUST be covered
- Use standard library `testing` package or `testify` for assertions

**Rationale**: Tests are documentation and safety nets. A CLI without tests is a CLI that breaks silently.

### III. CLI Design

The CLI MUST follow Unix conventions:
- Exit code 0 for success, non-zero for errors
- Errors written to stderr, output to stdout
- Support `--help` and `--version` flags
- Use clear, descriptive command and flag names
- Provide actionable error messages with context

**Rationale**: Users expect predictable CLI behavior. Consistent conventions enable scripting and automation.

## Technical Stack

- **Language**: Go 1.25
- **Testing**: `go test` with standard library or `testify`
- **Linting**: `go vet`, `staticcheck`, `gofmt`
- **Build**: `go build` with module support

## Governance

- This constitution supersedes conflicting practices
- Amendments require version bump and documented rationale
- All code changes MUST pass constitution checks before merge
- Versioning follows SemVer: MAJOR (breaking), MINOR (features), PATCH (fixes)

**Version**: 1.0.0 | **Ratified**: 2026-01-28 | **Last Amended**: 2026-01-28
