CREATE TABLE mas_ammar_k_wallet_member_saldo
(
    id_member   VARCHAR(30) PRIMARY KEY,
    last_saldo  NUMERIC(25, 2),
    last_update TIMESTAMP,
    user_id     VARCHAR(5),
    comm_code   VARCHAR(35)
);
