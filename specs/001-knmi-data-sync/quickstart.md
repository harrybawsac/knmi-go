# Quickstart: KNMI Weather Data Sync CLI

**Date**: 2026-01-28  
**Feature**: 001-knmi-data-sync

## Prerequisites

- Go 1.25 or later
- PostgreSQL 14+ (running locally or accessible)
- Internet access (to download KNMI data)

## Installation

```bash
# Clone the repository
git clone https://github.com/your-org/knmi-go.git
cd knmi-go

# Build the CLI
go build -o knmi ./cmd/knmi

# Verify installation
./knmi --version
```

## Database Setup

```bash
# Create the database
createdb knmi

# Set the connection string
export DATABASE_URL="postgres://localhost/knmi?sslmode=disable"
```

## Usage

### 1. Run Migrations

Apply the database schema:

```bash
./knmi migrate
```

Expected output:
```
Applied 1 migration:
  001_create_tables.sql
```

### 2. Sync Weather Data

Download and import KNMI weather data:

```bash
./knmi sync
```

Expected output:
```
Synced 45612 new records (45612 total)
```

### 3. Verify Data

Check the data in PostgreSQL:

```bash
psql $DATABASE_URL -c "SELECT date, tg/10.0 AS temp_c FROM weather_records ORDER BY date DESC LIMIT 5;"
```

## Common Tasks

### Run with Verbose Output

```bash
./knmi sync --verbose
```

### Use Custom Database

```bash
./knmi migrate --database-url "postgres://user:pass@remote-host:5432/weather"
```

### Override Data Source URL

```bash
./knmi sync --url "https://example.com/custom-data.zip"
```

### Check Version

```bash
./knmi --version
# or
./knmi version
```

## Troubleshooting

### "no migrations applied" Error

Run migrations before syncing:
```bash
./knmi migrate
./knmi sync
```

### Connection Refused

Check that PostgreSQL is running and `DATABASE_URL` is correct:
```bash
pg_isready -d $DATABASE_URL
```

### Permission Denied

Ensure the database user has CREATE TABLE permissions:
```bash
psql $DATABASE_URL -c "GRANT ALL ON DATABASE knmi TO your_user;"
```

## Development

### Run Tests

```bash
go test ./...
```

### Run with Coverage

```bash
go test -cover ./...
```

### Lint Code

```bash
go vet ./...
staticcheck ./...
```

## Project Structure

```
knmi-go/
├── cmd/knmi/           # CLI entrypoint
├── internal/           # Internal packages
│   ├── cli/            # Command implementations
│   ├── migration/      # Migration runner
│   ├── fetch/          # HTTP download & extraction
│   ├── parser/         # CSV parsing
│   ├── db/             # Database operations
│   └── config/         # Configuration
├── migrations/         # SQL migration files
└── tests/              # Test files
```
