CREATE TABLE IF NOT EXISTS themes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) UNIQUE NOT NULL,
    description VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    default_config JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);