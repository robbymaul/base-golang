alter table channels
    add column if not exists aggregator_id bigint;