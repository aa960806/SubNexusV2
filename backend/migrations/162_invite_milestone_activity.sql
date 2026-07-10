-- Cumulative invite milestone (累计邀请里程碑 / 成长进度条) activity.
-- Gameplay is controlled by ACTIVITY_CONFIG (invite_milestone_enabled + tiers);
-- this only adds a default activity-center entry. Visibility is gated by
-- isActivityCenterItemEnabledByFeature, so while the feature switch is off the
-- card stays hidden — identical to before it existed.
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
    'invite-milestone',
    '邀请里程碑',
    '累计邀请达标解锁阶梯宝箱奖励',
    '累计邀请新用户达到各档人数即可开启对应宝箱，领取阶梯奖励，直接发放到账户余额。',
    'gift',
    '/invite-milestone',
    '去领取',
    'invite_milestone',
    TRUE,
    80,
    '{}'::jsonb
)
ON CONFLICT (slug) DO NOTHING;
