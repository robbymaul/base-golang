alter table payments
    add column if not exists aggregator_id bigint,
    add constraint fk_payments_aggregators foreign key (aggregator_id) references aggregators(id);