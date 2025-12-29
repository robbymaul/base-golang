-- Drop the index first
DROP INDEX IF EXISTS idx_admin_users_uuid;

-- Then drop the column
ALTER TABLE admin_users
    DROP COLUMN IF EXISTS uuid;