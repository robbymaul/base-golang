alter table virtual_account_k_wallet
    add constraint fk_virtual_account_k_wallet_k_wallet_id foreign key (k_wallet_id) references k_wallet(id) on delete cascade ;