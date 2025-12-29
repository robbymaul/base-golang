alter table k_wallet_transaction
    add constraint fk_k_wallet_transaction_k_wallet_id foreign key (k_wallet_id) references k_wallet(id) on delete set null ;