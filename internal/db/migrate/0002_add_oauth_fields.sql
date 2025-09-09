-- Add OAuth fields to users table
ALTER TABLE users ADD COLUMN name TEXT;
ALTER TABLE users ADD COLUMN provider TEXT;
ALTER TABLE users ADD COLUMN provider_id TEXT;
ALTER TABLE users ADD COLUMN avatar_url TEXT;

-- Create index for provider lookup
CREATE INDEX idx_users_provider_id ON users(provider, provider_id);