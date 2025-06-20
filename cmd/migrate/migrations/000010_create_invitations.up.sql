CREATE TABLE IF NOT EXISTS user_invitations (
    token bytea PRIMARY KEY,
    user_id BIGINT NOT NULL,
    CONSTRAINT fk_user_id_invitations FOREIGN KEY (user_id) REFERENCES users(id)
);