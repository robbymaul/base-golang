create table if not exists topup_transaction(
    id bigint generated always as identity primary key not null,
    k_wallet_id bigint not null,
    member_id varchar(50),
    channel_id bigint,
    aggregator varchar(15),
    merchant varchar(50),
    amount numeric(15,2) not null,
    fee_admin numeric(15,2),
    currency varchar(5) not null,
    symbol varchar(5) not null,
    reference_id varchar(100),
    status varchar(100) not null,
    completed_at timestamp without time zone default current_timestamp,
    description text,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone default current_timestamp,
    deleted_at timestamp with time zone default null
)