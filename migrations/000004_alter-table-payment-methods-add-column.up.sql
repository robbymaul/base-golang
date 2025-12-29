alter table payment_methods
    add column if not exists payment_method   varchar(50),
    add column if not exists transaction_type varchar(50)