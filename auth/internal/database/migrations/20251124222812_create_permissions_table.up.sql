CREATE TABLE IF NOT EXISTS permissions(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(15) UNIQUE NOT NULL
);