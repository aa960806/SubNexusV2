-- Check-in anti-abuse: support freezing unpaid users' check-in rewards until
-- their first real payment. Additive column on activity_reward_logs, defaulting
-- to FALSE so every existing/other reward grant is unaffected.
--   frozen = TRUE  → the check-in day is recorded (keeps streak & once-per-day
--                    intact) but the amount is NOT yet credited to balance.
--   frozen = FALSE → normal granted reward (existing behavior for all sources).
ALTER TABLE activity_reward_logs
    ADD COLUMN IF NOT EXISTS frozen BOOLEAN NOT NULL DEFAULT FALSE;

-- Speeds up "this user's outstanding frozen check-in rewards" lookups/settlement.
CREATE INDEX IF NOT EXISTS idx_activity_reward_logs_frozen
    ON activity_reward_logs(user_id, source)
    WHERE frozen = TRUE;
