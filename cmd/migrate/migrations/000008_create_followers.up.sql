CREATE TABLE IF NOT EXISTS followers (

    user_id BIGINT NOT NULL,
    follower_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (user_id, follower_id),
    
    CONSTRAINT fk_user_id_followers
        FOREIGN KEY (user_id)
            REFERENCES users (id)
            ON DELETE CASCADE,
    CONSTRAINT fk_follower_id_followers
        FOREIGN KEY (follower_id)
            REFERENCES users (id)
            ON DELETE CASCADE

);