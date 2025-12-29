alter table platforms
    add column if not exists notification_url varchar(255);