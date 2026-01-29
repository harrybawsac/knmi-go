# CLI Interface Contract: knmi

**Version**: 1.0.0  
**Date**: 2026-01-28

## Global Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--help` | `-h` | bool | false | Show help message |
| `--version` | `-v` | bool | false | Show version |
| `--verbose` | | bool | false | Enable detailed progress logging |
| `--database-url` | | string | `$DATABASE_URL` | PostgreSQL connection string |

## Commands

### `knmi migrate`

Apply pending database migrations from the migrations directory.

**Usage**: `knmi migrate [flags]`

**Flags**:
| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--migrations-dir` | string | `./migrations` | Path to migrations directory |

**Exit Codes**:
| Code | Meaning |
|------|---------|
| 0 | All migrations applied successfully |
| 1 | Migration failed (invalid SQL, connection error) |

**Output (stdout)**:
```
Applied 3 migrations:
  001_create_tables.sql
  002_add_indexes.sql
  003_add_constraints.sql
```

**Output (--verbose)**:
```
Connecting to database...
Found 3 pending migrations
Applying 001_create_tables.sql...
  Creating table: weather_records
  Creating table: migrations
Applied 001_create_tables.sql (12ms)
...
```

**Errors (stderr)**:
```
Error: migration 002_add_indexes.sql failed: syntax error at position 42
```

---

### `knmi sync`

Download KNMI weather data and sync to database.

**Usage**: `knmi sync [flags]`

**Flags**:
| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--url` | string | KNMI default URL | Override data source URL |

**Exit Codes**:
| Code | Meaning |
|------|---------|
| 0 | Sync completed successfully |
| 1 | Sync failed (network, parse, database error) |

**Output (stdout)**:
```
Synced 125 new records (45612 total)
```

**Output (--verbose)**:
```
Downloading from https://cdn.knmi.nl/...
Downloaded 2.1 MB
Extracting archive...
Parsing CSV (45737 rows)...
Checking existing records...
Inserting 125 new records...
Synced 125 new records (45612 total)
```

**Errors (stderr)**:
```
Error: failed to download data: connection refused
Error: failed to parse CSV: unexpected column count at row 1234
Error: database error: connection refused
Error: no migrations applied. Run 'knmi migrate' first.
```

---

### `knmi version`

Display version information (alternative to `--version` flag).

**Usage**: `knmi version`

**Output (stdout)**:
```
knmi version 1.0.0
```

---

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://user:pass@localhost:5432/knmi?sslmode=disable` |
| `KNMI_MIGRATIONS_DIR` | Default migrations directory | `/app/migrations` |
| `KNMI_DATA_URL` | Default KNMI data URL | `https://cdn.knmi.nl/.../etmgeg_260.zip` |

## Connection String Format

```
postgres://[user[:password]@][host][:port][/database][?param=value]
```

**Examples**:
```
postgres://localhost/knmi
postgres://knmi:secret@db.example.com:5432/weather?sslmode=require
```

## Error Message Format

All errors follow the pattern:
```
Error: <action>: <specific problem>
```

Examples:
- `Error: connecting to database: connection refused`
- `Error: applying migration 002_indexes.sql: relation "weather_records" does not exist`
- `Error: parsing CSV row 1234: invalid integer value "abc"`
