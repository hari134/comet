CREATE DATABASE comet;

--bun:split

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,                     -- Unique internal user ID
    email VARCHAR(255) UNIQUE NOT NULL,        -- User's email (unique across providers)
    name VARCHAR(255),                         -- Full name of the user
    profile_picture_url TEXT,                  -- Optional profile picture URL
    registration_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Registration timestamp
    last_login TIMESTAMP,                      -- Last login timestamp
    CONSTRAINT email_format CHECK (email LIKE '%_@__%.__%') -- Ensures valid email format
);

--bun:split

CREATE TABLE IF NOT EXISTS auth_providers (
    id SERIAL PRIMARY KEY,                     -- Unique provider ID
    name VARCHAR(50) NOT NULL UNIQUE           -- Name of the provider (e.g., 'google', 'facebook', 'email')
);

--bun:split

CREATE TABLE IF NOT EXISTS user_auth_providers (
    id SERIAL PRIMARY KEY,                     -- Unique record ID
    user_id INT REFERENCES users(id) ON DELETE CASCADE,  -- Foreign key to users table
    provider_id INT REFERENCES auth_providers(id) ON DELETE CASCADE,  -- Foreign key to auth_providers
    provider_auth_id INT NOT NULL,             -- ID referencing provider-specific table (e.g., google_auth, facebook_auth)
    UNIQUE (user_id, provider_id),             -- Ensure each user can have only one record per provider
    UNIQUE (provider_id, provider_auth_id)     -- Ensure each provider-specific auth ID is unique
);

--bun:split

CREATE TABLE IF NOT EXISTS google_auth (
    id SERIAL PRIMARY KEY,                   -- Unique Google auth record ID
    google_id VARCHAR(255) UNIQUE NOT NULL   -- Google User ID (sub)
);
