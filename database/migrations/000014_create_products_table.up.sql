CREATE TYPE product_status AS ENUM ('draft', 'active', 'archived', 'out_of_stock');
/*
 We are using JSONB for product_properties to allow each tenant
 to define custom attributes for their products without altering
 the database schema. This enables flexibility in a multi-tenant
 environment, where different stores may have different product
 characteristics (e.g., color, size, material, voltage, capacity, etc.).
 JSONB also supports indexing and querying, so we can efficiently
 filter or search on specific attributes while keeping the schema
 stable and maintainable.
 */
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    price_cents BIGINT NOT NULL CHECK (price_cents >= 0),
    currency CHAR(3) NOT NULL DEFAULT 'KES',
    sku VARCHAR(100) UNIQUE NOT NULL,
    stock_quantity INT DEFAULT 0,
    product_properties JSONB NOT NULL DEFAULT '{}',
    product_status product_status NOT NULL DEFAULT 'draft',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
-- Index for fast tenant queries
CREATE INDEX idx_products_tenant ON products(tenant_id);