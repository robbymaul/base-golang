alter table platform_configurations
    drop column if exists config_name,
    drop column if exists config_json