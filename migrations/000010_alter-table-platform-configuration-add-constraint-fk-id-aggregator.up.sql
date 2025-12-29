alter table platform_configurations
    add column if not exists aggregator_id bigint,
    add constraint fk_id_aggregator foreign key (aggregator_id) references aggregators (id);