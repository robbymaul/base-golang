CREATE TABLE mas_ammar_k_wallet_member
(
    rec_id    SERIAL PRIMARY KEY,
    id_member VARCHAR(90)  NULL,
    bank_code VARCHAR(3)   NULL,
    number_va VARCHAR(50)  NULL,
    nama      VARCHAR(750) NULL,
    hp        VARCHAR(90)  NULL,
    email     VARCHAR(450) NULL,
    date_add  TIMESTAMP    NULL,
    date_upd  TIMESTAMP    NULL,
    status    SMALLINT     NULL,
    user_id   INTEGER      NULL
);
