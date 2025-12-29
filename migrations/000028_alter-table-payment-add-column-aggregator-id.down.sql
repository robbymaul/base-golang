alter table payments
    drop constraint if exists fk_payments_aggregators ,
    drop column if exists aggregator_id;
