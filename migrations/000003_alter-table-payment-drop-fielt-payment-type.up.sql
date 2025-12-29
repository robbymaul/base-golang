alter table payments
    drop constraint if exists payments_payment_type_id_fkey;

alter table payments
    drop column if exists payment_type_id