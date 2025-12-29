alter table payment_methods
    add column if not exists bank_name varchar(50);