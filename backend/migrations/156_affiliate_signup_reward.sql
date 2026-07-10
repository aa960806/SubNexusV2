ALTER TABLE user_affiliate_ledger
    ADD COLUMN IF NOT EXISTS source_ip VARCHAR(64) NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_user_affiliate_ledger_signup_reward_ip
    ON user_affiliate_ledger(action, source_ip, created_at)
    WHERE action = 'signup_bonus_inviter' AND source_ip <> '';

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_affiliate_signup_reward_inviter_once
    ON user_affiliate_ledger(source_user_id, action)
    WHERE action = 'signup_bonus_inviter' AND source_user_id IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_affiliate_signup_reward_invitee_once
    ON user_affiliate_ledger(source_user_id, action)
    WHERE action = 'signup_bonus_invitee' AND source_user_id IS NOT NULL;

COMMENT ON COLUMN user_affiliate_ledger.source_ip IS 'Source IP captured for affiliate signup reward anti-abuse checks.';
