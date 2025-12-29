alter table payment_methods
    drop column if exists payment_method,
    drop column if exists transaction_type