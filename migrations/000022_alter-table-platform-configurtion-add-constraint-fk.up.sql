alter table platform_configuration
    add constraint fk_platform_configuration_configuration_id foreign key (configuration_id) references configurations (id) on delete cascade,
    add constraint fk_platform_configuration_platform_id foreign key (platform_id) references platforms (id) on delete cascade;