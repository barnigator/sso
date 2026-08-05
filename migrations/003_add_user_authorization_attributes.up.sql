ALTER TABLE users
    ADD COLUMN role TEXT NOT NULL DEFAULT 'user',
    ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE;

ALTER TABLE users
ADD CONSTRAINT user_role_check
    CHECK ( role IN ('user', 'seller', 'admin') );

UPDATE users AS u
SET role = 'admin'
WHERE EXISTS (
    SELECT 1
    FROM admins AS a
    WHERE a.user_id = u.id
);