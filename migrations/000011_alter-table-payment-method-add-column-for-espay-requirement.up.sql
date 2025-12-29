alter table payment_methods
    add column if not exists is_espay boolean,
    add column if not exists product_name varchar(50),
    add column if not exists product_code varchar(50),
    add column if not exists bank_code varchar(50);