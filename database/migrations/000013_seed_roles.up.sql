INSERT INTO roles (name, slug, description, created_at)
VALUES('owner', 'owner', 'store owner', NOW()),
    (
        'system admin',
        'system-admin',
        'manages the entire system',
        NOW()
    );