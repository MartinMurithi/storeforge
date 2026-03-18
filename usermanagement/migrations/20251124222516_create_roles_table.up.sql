CREATE TABLE IF NOT EXISTS roles(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(15) UNIQUE NOT NULL,
    slug VARCHAR(15) UNIQUE NOT NULL,
    description TEXT,
    is_system BOOLEAN DEFAULT false
);