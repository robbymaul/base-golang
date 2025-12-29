CREATE TABLE mas_ammar_k_wallet_gen_va
(
    rec_id    SERIAL PRIMARY KEY,
    member_id VARCHAR(30) NULL ,
    user_id   VARCHAR(11) NULL ,
    date_add  TIMESTAMP NULL ,
    date_upd  TIMESTAMP NULL ,
    active    INTEGER NULL
);
