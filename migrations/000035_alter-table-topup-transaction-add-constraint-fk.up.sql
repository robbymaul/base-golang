alter table topup_transaction
    add constraint fk_topup_transaction_k_wallet_id foreign key (k_wallet_id) references k_wallet(id) on delete set null ,
    add constraint fk_topup_transaction_channel_id foreign key (channel_id) references channels(id) on delete set null ;