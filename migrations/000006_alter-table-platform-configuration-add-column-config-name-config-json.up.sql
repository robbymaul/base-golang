alter table platform_configurations
    add column if not exists config_name varchar(50),
    add column if not exists config_json json