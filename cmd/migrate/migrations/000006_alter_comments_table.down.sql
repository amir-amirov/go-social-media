ALTER TABLE comments
DROP CONSTRAINT IF EXISTS fk_user_id_comments,
DROP CONSTRAINT IF EXISTS fk_post_id_comments;