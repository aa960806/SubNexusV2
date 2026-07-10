-- Migration 174 keeps legacy activity data for database compatibility while disabling the four
-- removed gameplay entries so upgraded installations do not advertise routes
-- that are intentionally absent from SubNexusV2.
UPDATE activity_center_items
SET enabled = FALSE,
    updated_at = NOW()
WHERE slug IN (
    'daily-spin',
    'invite-lottery',
    'recharge-wheel',
    'invite-milestone'
)
  AND enabled = TRUE;
