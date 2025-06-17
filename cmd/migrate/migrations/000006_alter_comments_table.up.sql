ALTER TABLE comments
ADD CONSTRAINT fk_user_id_comments FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
ADD CONSTRAINT fk_post_id_comments FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE;