create table platform_channel
(
    id          bigint generated always as identity primary key not null,
    platform_id bigint                                          not null,
    channel_id  bigint                                          not null
);

create index idx_platform_channel_platform_id on platform_channel (platform_id);
create index idx_platform_channel_channel_id on platform_channel (channel_id);
