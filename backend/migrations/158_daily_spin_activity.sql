-- Daily usage threshold wheel activity.
-- The gameplay is controlled by ACTIVITY_CONFIG; this only adds a default
-- activity-center entry that admins can edit or disable like any other entry.
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
    'daily-spin',
    '每日消耗转盘',
    '今日消耗达标后抽取额度奖励',
    '每日消耗达到配置门槛后可参与一次转盘抽奖，次日重置。',
    'sparkles',
    '/daily-spin',
    '去抽奖',
    'daily_spin',
    TRUE,
    50,
    '{}'::jsonb
)
ON CONFLICT (slug) DO NOTHING;
