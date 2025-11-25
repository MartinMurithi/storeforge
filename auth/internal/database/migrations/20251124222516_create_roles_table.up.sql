CREATE TABLE IF NOT EXISTS roles(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(15) UNIQUE NOT NULL CHECK(name IN('admin', 'editor', 'member', 'viewer')),
    description TEXT
);