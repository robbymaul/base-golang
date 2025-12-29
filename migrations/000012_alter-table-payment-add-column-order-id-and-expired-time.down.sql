alter table payments
    drop column if exists order_id,
    drop column if exists expired_time;