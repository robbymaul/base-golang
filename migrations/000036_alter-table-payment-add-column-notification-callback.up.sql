alter table payments
    add column if not exists notification_callback boolean default null null ;