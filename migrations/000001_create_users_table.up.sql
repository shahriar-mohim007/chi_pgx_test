CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),-- UUID for each user
    name VARCHAR(100) NOT NULL,            -- User's full name
    email VARCHAR(255) NOT NULL UNIQUE,    -- User's email (must be unique)
    password VARCHAR(255) NOT NULL,        -- Hashed password for security
    is_active BOOLEAN DEFAULT FALSE,       -- Indicates if the user is activated
    created_at TIMESTAMP DEFAULT NOW(),    -- Timestamp of when the user was created
    updated_at TIMESTAMP DEFAULT NOW()     -- Timestamp of the last update
);

