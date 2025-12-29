alter table aggregators
    add column if not exists currency varchar(4) default 'IDR';