alter table platform_channel
    add constraint fk_platform_channel_platform_id
        foreign key (platform_id) references platforms (id) on delete cascade,
    add constraint fk_platform_channel_channel_id
        foreign key (channel_id) references payment_methods (id) on delete cascade;
