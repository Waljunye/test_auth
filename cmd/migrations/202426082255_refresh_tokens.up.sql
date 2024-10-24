-- migrate Up
CREATE TABLE IF NOT EXISTS refresh_tokens(
    user_uid uuid REFERENCES users(uuid) not null,
    token text not null,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    expires_at TIMESTAMP WITHOUT TIME ZONE
)