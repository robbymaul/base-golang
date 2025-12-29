alter table admin_users
    add column if not exists uuid uuid;

create index if not exists idx_admin_users_uuid on admin_users(uuid);