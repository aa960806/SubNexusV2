-- Recharge-reward double wheel (充值奖励双层转盘) activity.
-- Gameplay is controlled by ACTIVITY_CONFIG (recharge_wheel_enabled, threshold,
-- amounts + multipliers with probabilities); this only adds a default
-- activity-center entry. Visibility is gated by isActivityCenterItemEnabledByFeature,
-- so while the feature switch is off the card stays hidden — identical to before it existed.
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
    'recharge-wheel',
    '充值大转盘',
    '累计充值达标即可参与双层转盘抽奖',
    '累计充值每满配置门槛即赠送一次双层转盘机会，内层金额 × 外层倍数即为最终奖励，直接发放到账户余额。',
    'gift',
    '/recharge-wheel',
    '去抽奖',
    'recharge_wheel',
    TRUE,
    70,
    '{}'::jsonb
)
ON CONFLICT (slug) DO NOTHING;
