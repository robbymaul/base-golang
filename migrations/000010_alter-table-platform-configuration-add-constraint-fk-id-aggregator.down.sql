alter table platform_configurations
    drop constraint if exists fk_id_aggregator ,
    drop column if exists aggregator_id ;