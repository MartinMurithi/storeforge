CREATE TABLE IF NOT tenant_settings (
    tenant_id UUID PRIMARY KEY REFERENCES tenants(id) ON DELETE CASCADE,
    theme_id UUID NOT NULL REFERENCES themes(id),
    config JSONB NOT NULL DEFAULT '{}',
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_tenants_slug ON tenants(slug);