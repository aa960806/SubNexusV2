-- Activity rewards: leaderboard reward grants and daily check-ins.
CREATE TABLE IF NOT EXISTS activity_reward_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    source VARCHAR(32) NOT NULL,
    period VARCHAR(32) NOT NULL DEFAULT '',
    rank INTEGER NOT NULL DEFAULT 0,
    amount DECIMAL(20,8) NOT NULL DEFAULT 0,
    note TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT activity_reward_logs_positive_amount CHECK (amount > 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_activity_reward_logs_source_period_user
    ON activity_reward_logs(source, period, user_id);

CREATE INDEX IF NOT EXISTS idx_activity_reward_logs_user_created
    ON activity_reward_logs(user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_activity_reward_logs_source_created
    ON activity_reward_logs(source, created_at DESC);
