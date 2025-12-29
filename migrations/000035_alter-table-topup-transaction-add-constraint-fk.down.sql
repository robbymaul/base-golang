alter table topup_transaction
    drop constraint if exists fk_topup_transaction_k_wallet_id,
    drop constraint if exists fk_topup_transaction_channel_id;