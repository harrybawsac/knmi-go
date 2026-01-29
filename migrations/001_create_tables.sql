-- Migration: 001_create_tables.sql
-- Creates the initial database schema for KNMI weather data.

-- Table: migrations
-- Tracks which migrations have been applied to the database.
CREATE TABLE IF NOT EXISTS migrations (
    id SERIAL PRIMARY KEY,
    version INTEGER NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    applied_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Table: weather_records
-- Stores daily weather observations from KNMI stations.
CREATE TABLE IF NOT EXISTS weather_records (
    id SERIAL PRIMARY KEY,
    station_id INTEGER NOT NULL,
    date DATE NOT NULL,
    
    -- Wind measurements
    ddvec INTEGER,          -- Vector mean wind direction (degrees)
    fhvec INTEGER,          -- Vector mean windspeed (0.1 m/s)
    fg INTEGER,             -- Daily mean windspeed (0.1 m/s)
    fhx INTEGER,            -- Max hourly windspeed (0.1 m/s)
    fhxh INTEGER,           -- Hour of FHX
    fhn INTEGER,            -- Min hourly windspeed (0.1 m/s)
    fhnh INTEGER,           -- Hour of FHN
    fxx INTEGER,            -- Max wind gust (0.1 m/s)
    fxxh INTEGER,           -- Hour of FXX
    
    -- Temperature measurements
    tg INTEGER,             -- Daily mean temp (0.1 °C)
    tn INTEGER,             -- Min temp (0.1 °C)
    tnh INTEGER,            -- Hour of TN
    tx INTEGER,             -- Max temp (0.1 °C)
    txh INTEGER,            -- Hour of TX
    t10n INTEGER,           -- Min temp at 10cm (0.1 °C)
    t10nh INTEGER,          -- Period of T10N
    
    -- Sunshine and radiation
    sq INTEGER,             -- Sunshine duration (0.1 hour)
    sp INTEGER,             -- Sunshine percentage
    q INTEGER,              -- Global radiation (J/cm²)
    
    -- Precipitation
    dr INTEGER,             -- Precipitation duration (0.1 hour)
    rh INTEGER,             -- Daily precipitation (0.1 mm)
    rhx INTEGER,            -- Max hourly precipitation (0.1 mm)
    rhxh INTEGER,           -- Hour of RHX
    
    -- Pressure
    pg INTEGER,             -- Mean sea level pressure (0.1 hPa)
    px INTEGER,             -- Max pressure (0.1 hPa)
    pxh INTEGER,            -- Hour of PX
    pn INTEGER,             -- Min pressure (0.1 hPa)
    pnh INTEGER,            -- Hour of PN
    
    -- Visibility
    vvn INTEGER,            -- Min visibility (coded)
    vvnh INTEGER,           -- Hour of VVN
    vvx INTEGER,            -- Max visibility (coded)
    vvxh INTEGER,           -- Hour of VVX
    
    -- Cloud cover and humidity
    ng INTEGER,             -- Mean cloud cover (octants)
    ug INTEGER,             -- Mean relative humidity (%)
    ux INTEGER,             -- Max relative humidity (%)
    uxh INTEGER,            -- Hour of UX
    un INTEGER,             -- Min relative humidity (%)
    unh INTEGER,            -- Hour of UN
    
    -- Evapotranspiration
    ev24 INTEGER,           -- Potential evapotranspiration (0.1 mm)
    
    -- Metadata
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT weather_records_station_date_unique UNIQUE (station_id, date)
);

-- Index for fast lookups by station and date
CREATE INDEX IF NOT EXISTS idx_weather_records_station_date 
    ON weather_records (station_id, date);
