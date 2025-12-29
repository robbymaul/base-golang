create table if not exists mas_dion_va_cust_pay_bal (
    id int primary key not null ,
    trcd varchar(20) null ,
    trdt timestamp null ,
    novac varchar(20) null ,
    dfno varchar(100) null ,
    fullnm varchar(100) null ,
    "type" char(1) default '0',
    refno varchar(100) null,
    amount numeric(15,2) null ,
    status char(1) default 'O',
    custtype varchar(1) default 'O',
    description varchar(150) null,
    remarks varchar(100) null
)