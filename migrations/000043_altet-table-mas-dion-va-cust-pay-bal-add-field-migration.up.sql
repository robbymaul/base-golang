alter table mas_dion_va_cust_pay_bal
    add column if not exists migration bool default false;