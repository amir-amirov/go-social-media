-- This may not work in production
-- ALTER TABLE IF EXISTS users 
-- ADD COLUMN IF NOT EXISTS role_id BIGINT DEFAULT 1,
-- ADD CONSTRAINT fk_role_id_users FOREIGN KEY (role_id) REFERENCES roles(id);

-- So I did:

-- Add Column First: no default, no constraint
ALTER TABLE IF EXISTS users 
ADD COLUMN IF NOT EXISTS role_id BIGINT;

-- Update with Dynamic Role ID
UPDATE users
SET role_id = (
    SELECT id FROM roles WHERE name = 'user'
);

-- Add Foreign Key Constraint AFTER Data Clean
ALTER TABLE users
ADD CONSTRAINT fk_role_id_users FOREIGN KEY (role_id) REFERENCES roles(id);
-- Now FK applies only after valid data exists.

-- Make Column Non-Nullable
ALTER TABLE users
ALTER COLUMN role_id SET NOT NULL;
