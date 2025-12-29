alter table payment_methods
    drop constraint if exists fk_id_platform ,
    drop column if exists platform_id;
