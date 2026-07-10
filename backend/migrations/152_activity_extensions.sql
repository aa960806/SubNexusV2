-- Activity extension switches, broadcast marquee messages, and first recharge gift tracking.
CREATE TABLE IF NOT EXISTS activity_broadcasts (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(120) NOT NULL DEFAULT '',
    content TEXT NOT NULL,
    source VARCHAR(32) NOT NULL DEFAULT 'admin',
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    priority INTEGER NOT NULL DEFAULT 0,
    start_at TIMESTAMPTZ,
    end_at TIMESTAMPTZ,
    created_by BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_activity_broadcasts_active
    ON activity_broadcasts(enabled, priority DESC, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_activity_broadcasts_window
    ON activity_broadcasts(start_at, end_at);

CREATE TABLE IF NOT EXISTS first_recharge_gift_purchases (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    order_id BIGINT REFERENCES payment_orders(id) ON DELETE SET NULL,
    price DECIMAL(20,2) NOT NULL DEFAULT 0,
    credited_amount DECIMAL(20,2) NOT NULL DEFAULT 0,
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_first_recharge_gift_purchases_user
    ON first_recharge_gift_purchases(user_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_first_recharge_gift_purchases_order
    ON first_recharge_gift_purchases(order_id)
    WHERE order_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_first_recharge_gift_purchases_status
    ON first_recharge_gift_purchases(status, created_at DESC);
