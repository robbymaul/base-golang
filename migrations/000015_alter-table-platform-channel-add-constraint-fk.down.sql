alter table platform_channel
    drop constraint IF EXISTS fk_platform_channel_platform_id,
    drop constraint IF EXISTS fk_platform_channel_channel_id;
