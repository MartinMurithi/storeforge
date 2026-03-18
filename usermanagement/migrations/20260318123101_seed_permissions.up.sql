INSERT INTO permissions (slug, category, description)
VALUES -- Product Management
    (
        'products:read',
        'products',
        'Can view product list and details'
    ),
    (
        'products:write',
        'products',
        'Can create and edit products'
    ),
    (
        'products:delete',
        'products',
        'Can permanently remove products'
    ),
    -- Order Management
    (
        'orders:read',
        'orders',
        'Can view customer orders'
    ),
    (
        'orders:write',
        'orders',
        'Can update order status and fulfillment'
    ),
    -- Tenant/Store Settings
    (
        'settings:read',
        'settings',
        'Can view store configuration'
    ),
    (
        'settings:write',
        'settings',
        'Can update store theme and metadata'
    ),
    -- User/Team Management
    (
        'users:invite',
        'users',
        'Can invite new members to the tenant'
    ),
    (
        'users:manage',
        'users',
        'Can change roles or remove members'
    ),
    -- Super Admin Only (Platform Level)
    (
        'tenants:create',
        'super-admin',
        'Can create new store instances'
    ),
    (
        'tenants:suspend',
        'super-admin',
        'Can disable a store instance'
    ) ON CONFLICT(slug) DO NOTHING;