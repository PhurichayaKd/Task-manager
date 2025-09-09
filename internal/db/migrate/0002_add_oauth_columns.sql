-- internal/db/migrate/0002_add_oauth_columns.sql
ALTER TABLE users ADD COLUMN name TEXT;
ALTER TABLE users ADD COLUMN provider VARCHAR(32);
ALTER TABLE users ADD COLUMN provider_id VARCHAR(191);
ALTER TABLE users ADD COLUMN avatar_url TEXT;

CREATE INDEX IF NOT EXISTS idx_users_provider_id ON users(provider, provider_id);
