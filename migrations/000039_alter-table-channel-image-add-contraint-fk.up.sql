alter table channel_images
    add constraint fk_channel_images_channel_id foreign key (channel_id) references channels (id);