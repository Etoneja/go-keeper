CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS secrets (
    id VARCHAR(255) PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    last_modified TIMESTAMP NOT NULL,
    hash VARCHAR(255) NOT NULL,
    data BYTEA NOT NULL
);

ALTER TABLE secrets ADD CONSTRAINT secrets_id_user_id_key UNIQUE (id, user_id);
CREATE INDEX IF NOT EXISTS idx_secrets_user_id ON secrets(user_id);
