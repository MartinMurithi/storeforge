CREATE TABLE IF NOT EXISTS tenants(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_name VARCHAR(25) UNIQUE NOT NULL,
    business_type VARCHAR(25) UNIQUE NOT NULL,
    slug VARCHAR(25) UNIQUE NOT NULL,
    sub_domain VARCHAR(25) UNIQUE NOT NULL,
    status VARCHAR(25) NOT NULL DEFAULT 'provisioning' CHECK (
        status IN(
            'provisioning',
            'active',
            'suspended',
            'pending_deletion',
            'deleted'
        )
    ),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NOW()
);