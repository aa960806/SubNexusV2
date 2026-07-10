-- Activity check-in IP constraint: record the client IP on each reward grant so the
-- service can enforce the admin-toggleable "one check-in per IP per day" rule.
ALTER TABLE activity_reward_logs
    ADD COLUMN IF NOT EXISTS ip VARCHAR(64) NOT NULL DEFAULT '';

-- Supports the per-period IP lookup performed before each check-in grant.
-- Enforcement lives in the service layer (gated by ACTIVITY_CONFIG.checkin_ip_limit),
-- so this index is intentionally non-unique.
CREATE INDEX IF NOT EXISTS idx_activity_reward_logs_source_period_ip
    ON activity_reward_logs(source, period, ip)
    WHERE ip <> '';
