-- Site uptime calendar.
-- Heartbeats are sampled by the application process. Calendar days can be
-- calculated automatically, while manual admin overrides are kept separate.
CREATE TABLE IF NOT EXISTS uptime_calendar_heartbeats (
    id BIGSERIAL PRIMARY KEY,
    bucket_at TIMESTAMPTZ NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_uptime_calendar_heartbeats_bucket
    ON uptime_calendar_heartbeats(bucket_at DESC);

CREATE TABLE IF NOT EXISTS uptime_calendar_days (
    day DATE PRIMARY KEY,
    status VARCHAR(12) NOT NULL DEFAULT 'green',
    source VARCHAR(12) NOT NULL DEFAULT 'auto',
    outage_minutes INTEGER NOT NULL DEFAULT 0,
    note TEXT NOT NULL DEFAULT '',
    manual_by BIGINT,
    manual_at TIMESTAMPTZ,
    calculated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_uptime_calendar_days_status CHECK (status IN ('green', 'red')),
    CONSTRAINT chk_uptime_calendar_days_source CHECK (source IN ('auto', 'manual'))
);

CREATE INDEX IF NOT EXISTS idx_uptime_calendar_days_day
    ON uptime_calendar_days(day DESC);

CREATE INDEX IF NOT EXISTS idx_uptime_calendar_days_status
    ON uptime_calendar_days(status, day DESC);
