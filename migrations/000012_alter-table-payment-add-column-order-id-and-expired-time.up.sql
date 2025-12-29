alter table payments
    add column if not exists order_id varchar(100),
    add column if not exists expired_time varchar(100);