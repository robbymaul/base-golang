create table if not exists platform_configuration
(
    id               bigint generated always as identity primary key not null,
    configuration_id bigint                                          not null,
    platform_id      bigint                                          not null
);

create index idx_platform_configuration_configuration_id on platform_configuration (configuration_id);
create index idx_platform_configuration_platform_id on platform_configuration (platform_id);
