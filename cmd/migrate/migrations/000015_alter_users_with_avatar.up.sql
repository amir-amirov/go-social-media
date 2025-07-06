ALTER TABLE IF EXISTS users 
ADD COLUMN IF NOT EXISTS avatar TEXT DEFAULT 'https://firebasestorage.googleapis.com/v0/b/auth-2c46a.appspot.com/o/user-profile-icon-profile-avatar-user-icon-male-icon-face-icon-profile-icon-free-png.webp?alt=media&token=a96e4694-e39e-4d96-80be-034d19eebb52';

UPDATE users
SET avatar = 'https://firebasestorage.googleapis.com/v0/b/auth-2c46a.appspot.com/o/user-profile-icon-profile-avatar-user-icon-male-icon-face-icon-profile-icon-free-png.webp?alt=media&token=a96e4694-e39e-4d96-80be-034d19eebb52'
WHERE avatar IS NULL;