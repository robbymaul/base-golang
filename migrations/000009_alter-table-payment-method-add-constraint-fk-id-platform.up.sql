alter table payment_methods
    add column if not exists platform_id bigint,
    add constraint fk_id_platform foreign key (platform_id) references platforms (id);