alter table payment_methods
    add constraint fk_id_platform foreign key (platform_id) references platforms (id) on delete set null ;