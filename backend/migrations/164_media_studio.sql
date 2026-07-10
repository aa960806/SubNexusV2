-- Media Studio (platform-hosted media generation) — image MVP.
-- Fully self-contained feature: these tables are independent from the account
-- pool, usage_logs and existing billing/stats. When the feature switch
-- media_studio_enabled is off, nothing here is read by user-facing endpoints.
--
-- Upstream base URLs and API keys live ONLY in these server-side tables; the
-- api_key_encrypted column stores an AES-256-GCM ciphertext (never plaintext),
-- and is never exposed to any user-facing endpoint or public settings payload.

-- Upstream channels (provider endpoint + encrypted key). Admin-only.
CREATE TABLE IF NOT EXISTS media_upstreams (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(120) NOT NULL,
    media_types JSONB NOT NULL DEFAULT '["image"]'::jsonb,
    protocol VARCHAR(32) NOT NULL DEFAULT 'openai_images',
    base_url TEXT NOT NULL DEFAULT '',
    api_key_encrypted TEXT NOT NULL DEFAULT '',
    status VARCHAR(16) NOT NULL DEFAULT 'enabled',
    timeout_seconds INTEGER NOT NULL DEFAULT 120,
    max_concurrency INTEGER NOT NULL DEFAULT 0,
    remark VARCHAR(255) NOT NULL DEFAULT '',
    extra JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_media_upstreams_status
    ON media_upstreams(status);

-- Platform-facing models. public_id is what the frontend/user sees; the real
-- upstream_model + upstream_id are server-side only and never leave the backend.
CREATE TABLE IF NOT EXISTS media_models (
    id BIGSERIAL PRIMARY KEY,
    public_id VARCHAR(80) NOT NULL UNIQUE,
    display_name VARCHAR(120) NOT NULL,
    media_type VARCHAR(16) NOT NULL DEFAULT 'image',
    upstream_id BIGINT NOT NULL REFERENCES media_upstreams(id) ON DELETE CASCADE,
    upstream_model VARCHAR(120) NOT NULL DEFAULT '',
    capabilities JSONB NOT NULL DEFAULT '["text_to_image"]'::jsonb,
    param_schema JSONB NOT NULL DEFAULT '{}'::jsonb,
    pricing JSONB NOT NULL DEFAULT '{}'::jsonb,
    tags JSONB NOT NULL DEFAULT '[]'::jsonb,
    sort_order INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(16) NOT NULL DEFAULT 'enabled',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_media_models_visible
    ON media_models(media_type, status, sort_order ASC, id ASC);

CREATE INDEX IF NOT EXISTS idx_media_models_upstream
    ON media_models(upstream_id);

-- Server-side spend/audit trail for media generation. Intentionally SEPARATE
-- from usage_logs (which has NOT NULL FKs to api_keys/accounts that media
-- generation does not have) so this feature never touches existing billing or
-- usage dashboards. Result images are stored client-side only, not here.
CREATE TABLE IF NOT EXISTS media_generation_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    media_type VARCHAR(16) NOT NULL DEFAULT 'image',
    mode VARCHAR(32) NOT NULL DEFAULT 'text_to_image',
    model_public_id VARCHAR(80) NOT NULL DEFAULT '',
    upstream_id BIGINT,
    upstream_model VARCHAR(120) NOT NULL DEFAULT '',
    prompt TEXT NOT NULL DEFAULT '',
    params JSONB NOT NULL DEFAULT '{}'::jsonb,
    image_count INTEGER NOT NULL DEFAULT 0,
    price DECIMAL(20, 8) NOT NULL DEFAULT 0,
    status VARCHAR(16) NOT NULL DEFAULT 'succeeded',
    error_message TEXT NOT NULL DEFAULT '',
    request_id VARCHAR(64) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_media_generation_logs_user_created
    ON media_generation_logs(user_id, created_at DESC);
