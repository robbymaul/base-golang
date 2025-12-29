CREATE TABLE IF NOT EXISTS mas_ammar_k_wallet_transaction
(
    id_trans         SERIAL PRIMARY KEY,
    id_member        VARCHAR(50) NULL ,
    link_id_member   VARCHAR(50) NULL ,
    order_id         VARCHAR(50) NULL ,
    comm_code        VARCHAR(25) NULL ,
    datetime         TIMESTAMP NULL ,
    valuedate        DATE NULL ,
    description      TEXT NULL ,
    reference_no     INTEGER NULL ,
    debit            INTEGER NULL ,
    credit           INTEGER NULL ,
    balance_orderid  INTEGER NULL ,
    balance_memberid INTEGER NULL ,
    balance_userid   INTEGER NULL ,
    balance_commcode INTEGER NULL ,
    balance_all      INTEGER NULL ,
    user_id          INTEGER NULL ,
    type_transaction VARCHAR(50) NULL ,
    is_prod          INTEGER NULL
);
