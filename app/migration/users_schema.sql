CREATE TABLE IF NOT EXISTS "users" 
(
    id UUID NOT NULL PRIMARY KEY DEFAULT (uuid_generate_v4()),
    name VARCHAR(100),
    email VARCHAR(255) NOT NULL UNIQUE,
    photo VARCHAR(300) NOT NULL DEFAULT 'default.png',
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);

