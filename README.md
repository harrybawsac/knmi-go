# knmi

A CLI tool to fetch weather data from KNMI (Royal Netherlands Meteorological Institute) and sync it to a PostgreSQL database.

## Features

- **Database Migrations**: Manage schema with SQL migration files
- **Incremental Sync**: Download KNMI weather data and insert only new records
- **Duplicate Prevention**: Uses `(station_id, date)` unique constraint to prevent duplicates
- **Configurable**: Override data source URL and database connection via flags or environment variables

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/harrybawsac/knmi-go.git
cd knmi-go

# Build the binary
go build -o bin/knmi ./cmd/knmi

# Install to your PATH (optional)
cp bin/knmi /usr/local/bin/
```

## Usage

### Set Up Database

1. Create a PostgreSQL database:

```bash
createdb knmi
```

2. Set the database URL:

```bash
export DATABASE_URL="postgres://user:password@localhost:5432/knmi?sslmode=disable"
```

3. Run migrations to create the schema:

```bash
knmi migrate
```

### Sync Weather Data

Download and sync the latest weather data from KNMI:

```bash
knmi sync
```

With verbose output:

```bash
knmi sync --verbose
```

### Commands

| Command | Description |
|---------|-------------|
| `knmi migrate` | Apply pending database migrations |
| `knmi sync` | Download and sync KNMI weather data |
| `knmi version` | Display version information |
| `knmi help` | Display help information |

### Global Flags

| Flag | Description |
|------|-------------|
| `--verbose` | Enable detailed progress logging |
| `--database-url` | PostgreSQL connection string (overrides DATABASE_URL) |
| `-h, --help` | Display help |

### Environment Variables

| Variable | Description |
|----------|-------------|
| `DATABASE_URL` | PostgreSQL connection string |
| `KNMI_DATA_URL` | Override default KNMI data URL |
| `KNMI_MIGRATIONS_DIR` | Path to migrations directory |

## Data Source

Weather data is fetched from KNMI station 260 (De Bilt):
- URL: https://cdn.knmi.nl/knmi/map/page/klimatologie/gegevens/daggegevens/etmgeg_260.zip
- Contains daily weather observations with 41 data columns

## Development

### Prerequisites

- Go 1.25+
- PostgreSQL 13+

### Build

```bash
make build
```

### Run Tests

```bash
# Unit tests
make test

# With coverage
go test -cover ./...
```

### Lint

```bash
make lint
```

## License

MIT
