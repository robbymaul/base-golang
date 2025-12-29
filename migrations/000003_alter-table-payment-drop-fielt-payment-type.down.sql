ALTER TABLE payments
    ADD COLUMN IF NOT EXISTS payment_type_id BIGINT;

ALTER TABLE payments
    ADD CONSTRAINT  payments_payment_type_id_fkey FOREIGN KEY (payment_type_id) REFERENCES payment_types(id);