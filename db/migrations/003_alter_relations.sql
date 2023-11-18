SET search_path TO pinspire;

ALTER TABLE profile ALTER COLUMN avatar 
SET DEFAULT 'https://pinspire.online:8081/upload/avatars/default-avatar.png';

UPDATE profile SET avatar = DEFAULT WHERE avatar = 'default-avatar.png';

ALTER TABLE profile ADD COLUMN about_me text;
