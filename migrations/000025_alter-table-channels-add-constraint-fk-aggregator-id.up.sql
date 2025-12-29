alter table channels
    add constraint fk_channels_aggregator_id foreign key (aggregator_id) references aggregators (id) on delete set null ;