# CLI Contract: Sync Command with Dry-Run Mode

**Feature**: 002-dry-run | **Date**: 2026-01-28

## Command Signature

```bash
knmi sync [flags]
```

## Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--dry-run` | `-n` | bool | false | Preview records without inserting into database |

## Behavior Matrix

| Mode | Download | Parse | Insert | Output |
|------|----------|-------|--------|--------|
| Normal (`knmi sync`) | ✅ | ✅ | ✅ | Summary of inserted records |
| Dry-run (`knmi sync --dry-run`) | ✅ | ✅ | ❌ | Preview table of last 10 new records |

## Output Format (Dry-Run Mode)

### Standard Output (stdout)

```
Dry-run mode: previewing records that would be inserted

DATE        STATION_ID     TG     TN     TX     FG     RH
2024-01-15         260    45     12     78     35    120
2024-01-15         270    43     10     75     32    115
2024-01-15         280    47     14     80     38    125
...

Total: 156 new records would be inserted (showing last 10)
```

### Output Fields

| Field | Source | Format | Null Display |
|-------|--------|--------|--------------|
| DATE | WeatherRecord.Date | YYYY-MM-DD | N/A (required) |
| STATION_ID | WeatherRecord.StationID | integer | N/A (required) |
| TG | WeatherRecord.TG | integer (0.1°C) | `-` |
| TN | WeatherRecord.TN | integer (0.1°C) | `-` |
| TX | WeatherRecord.TX | integer (0.1°C) | `-` |
| FG | WeatherRecord.FG | integer (0.1 m/s) | `-` |
| RH | WeatherRecord.RH | integer (0.1 mm) | `-` |

### Edge Case Outputs

**Fewer than 10 new records (e.g., 3 records):**
```
Dry-run mode: previewing records that would be inserted

DATE        STATION_ID     TG     TN     TX     FG     RH
2024-01-15         260    45     12     78     35    120
2024-01-15         270    43     10     75     32    115
2024-01-15         280    47     14     80     38    125

Total: 3 new records would be inserted (showing all)
```

**No new records:**
```
Dry-run mode: no new records to insert

All records from the KNMI file already exist in the database.
```

### Error Output (stderr)

**Database connection failure:**
```
Error: failed to connect to database: <error details>
```

**KNMI download failure:**
```
Error: failed to download KNMI data: <error details>
```

## Exit Codes

| Code | Condition |
|------|-----------|
| 0 | Success (with or without records to insert) |
| 1 | Error (DB connection, download, parse failure) |

## Examples

```bash
# Preview what would be synced
knmi sync --dry-run

# Short flag form
knmi sync -n

# Normal sync (unchanged behavior)
knmi sync
```

## Compatibility

- This is a **non-breaking change**
- Default behavior (`knmi sync`) remains unchanged
- Flag can be combined with any existing sync flags (if any added later)
