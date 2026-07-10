-- Invite-lottery (九宫格轮盘) activity.
-- The gameplay is controlled by ACTIVITY_CONFIG (invite_lottery_enabled +
-- invite_lottery_prizes); this only adds a default activity-center entry.
-- Visibility is gated by isActivityCenterItemEnabledByFeature, so while the
-- feature switch is off the card stays hidden — identical to before it existed.
INSERT INTO activity_center_items (
    slug,
    title,
    subtitle,
    description,
    icon,
    route_path,
    action_label,
    activity_type,
    enabled,
    sort_order,
    metadata
)
VALUES (
    'invite-lottery',
    '邀请抽奖',
    '每邀请一位新用户得一次抽奖机会',
    '每成功邀请一位新用户即赠送一次九宫格抽奖机会，抽中奖励直接发放到账户余额。',
    'gift',
    '/invite-lottery',
    '去抽奖',
    'invite_lottery',
    TRUE,
    60,
    '{}'::jsonb
)
ON CONFLICT (slug) DO NOTHING;
