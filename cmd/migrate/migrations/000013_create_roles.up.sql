CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    level INT NOT NULL DEFAULT 0,
    description TEXT
);

INSERT INTO roles (name, description, level)
VALUES ('user', 'A user can create, update and delete post and comments', 1),
    ('moderator', 'A moderator can update other users post', 2),
    ('admin', 'An admin can update and delete other users post', 3);