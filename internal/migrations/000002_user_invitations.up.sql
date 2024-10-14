CREATE TABLE IF NOT EXISTS user_invitations (
    user_id SERIAL PRIMARY KEY,
    token bytea NOT NULL,
    issue_at TIMESTAMP DEFAULT NOW(),
    expire_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)