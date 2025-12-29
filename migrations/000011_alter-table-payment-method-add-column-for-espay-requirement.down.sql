alter table payment_methods
    drop column if exists is_espay ,
    drop column if exists product_name,
    drop column if exists product_code ,
    drop column if exists bank_code ;