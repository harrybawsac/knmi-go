# Quickstart: Dry-Run Mode for Sync Command

**Feature**: 002-dry-run | **Date**: 2026-01-28

## Overview

The `--dry-run` flag allows you to preview which weather records would be inserted into the database without actually inserting them. This is useful for:

- Verifying data before committing to the database
- Testing the sync process without side effects
- Checking for new records before a scheduled sync

## Usage

### Preview Records

```bash
# Preview what would be synced (long form)
knmi sync --dry-run

# Preview what would be synced (short form)
knmi sync -n
```

### Example Output

```
Dry-run mode: previewing records that would be inserted

DATE        STATION_ID     TG     TN     TX     FG     RH
2024-01-15         260    45     12     78     35    120
2024-01-15         270    43     10     75     32    115
2024-01-15         280    47     14     80     38    125
2024-01-14         260    42     10     72     30    100
2024-01-14         270    40      8     70     28     95
2024-01-14         280    44     12     74     33    105
2024-01-13         260    38      5     68     25     80
2024-01-13         270    36      3     65     22     75
2024-01-13         280    40      7     70     28     85
2024-01-12         260    35      2     62     20     60

Total: 156 new records would be inserted (showing last 10)
```

### Understanding the Output

| Column | Description | Unit |
|--------|-------------|------|
| DATE | Measurement date | YYYY-MM-DD |
| STATION_ID | KNMI weather station ID | - |
| TG | Mean temperature | 0.1 °C |
| TN | Minimum temperature | 0.1 °C |
| TX | Maximum temperature | 0.1 °C |
| FG | Mean wind speed | 0.1 m/s |
| RH | Precipitation sum | 0.1 mm |

**Note**: A `-` in any field indicates the measurement was not available.

## Common Scenarios

### Check Before First Sync

```bash
# See all records that will be inserted on first sync
knmi sync --dry-run
```

### Verify Daily Updates

```bash
# After yesterday's sync, see today's new records
knmi sync -n
```

### No New Records

If all records already exist in your database:

```
Dry-run mode: no new records to insert

All records from the KNMI file already exist in the database.
```

## Requirements

- Database must be accessible (dry-run queries for existing records)
- Network access to KNMI servers (dry-run downloads the data file)
