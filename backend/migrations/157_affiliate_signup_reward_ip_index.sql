-- Widen the affiliate signup-reward IP index to cover BOTH reward actions.
--
-- The original index (migration 156) was partial on action = 'signup_bonus_inviter'
-- only. When the admin configures "reward invitee only" (inviter amount = 0), no
-- 'signup_bonus_inviter' rows are written, so the per-IP daily counter - which now
-- counts DISTINCT source_user_id across both actions - could not use the index and,
-- worse, the old query missed those rows entirely. This recreates the index over the
-- full reward-action set so the anti-abuse IP limit is enforceable in every config.
DROP INDEX IF EXISTS idx_user_affiliate_ledger_signup_reward_ip;

CREATE INDEX IF NOT EXISTS idx_user_affiliate_ledger_signup_reward_ip
    ON user_affiliate_ledger(source_ip, created_at, source_user_id)
    WHERE action IN ('signup_bonus_inviter', 'signup_bonus_invitee') AND source_ip <> '';
