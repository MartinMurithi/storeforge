INSERT INTO themes (
        id,
        name,
        description,
        is_active,
        default_config,
        created_at
    )
VALUES (
        '00000000-0000-0000-0000-000000000001',
        'Midnight Modern',
        'A sleek, dark-themed template optimized for electronics and high-end fashion.',
        true,
        '{
        "colors": {"primary": "#000000", "secondary": "#ffffff", "accent": "#ff4400"},
        "typography": {"font_family": "Inter, sans-serif", "base_size": "16px"},
        "layout": {"header": "sticky", "footer": "minimal"}
    }',
        NOW()
    ),
    (
        '00000000-0000-0000-0000-000000000002',
        'Artisan Bloom',
        'Soft palettes and elegant serif fonts, perfect for handmade crafts and florists.',
        true,
        '{
        "colors": {"primary": "#fdfaf6", "secondary": "#4a5d4e", "accent": "#d4a373"},
        "typography": {"font_family": "Playfair Display, serif", "base_size": "18px"},
        "layout": {"header": "centered", "footer": "extended"}
    }',
        NOW()
    );
