-- ============================================================================
-- PAYMENT SERVICE DATABASE SCHEMA
-- Microservice untuk menangani berbagai jenis pembayaran
-- ============================================================================

-- Create ENUM types at the beginning of your schema
-- CREATE TYPE  fee_type_enum AS ENUM ('fixed', 'percentage', 'none');
-- CREATE TYPE payment_status_enum AS ENUM ('pending', 'processing', 'success', 'failed', 'cancelled', 'expired');

-- Tabel untuk menyimpan informasi platform/aplikasi
CREATE TABLE platforms
(
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY NOT NULL,
    code        VARCHAR(50) UNIQUE                              NOT NULL, -- 'web_knet', 'web_sms', 'web_smartd', etc
    name        VARCHAR(100)                                    NOT NULL,
    description TEXT                                            NULL,
    api_key     VARCHAR(255) UNIQUE                             NOT NULL,
    secret_key  VARCHAR(255)                                    NOT NULL,
    is_active   BOOLEAN                  DEFAULT TRUE,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP WITH TIME ZONE                        NULL
);

-- Tabel untuk menyimpan jenis pembayaran
CREATE TABLE payment_types
(
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY NOT NULL,
    code        VARCHAR(50) UNIQUE                              NOT NULL, -- 'sales_order', 'topup_token', 'topup_wallet'
    name        VARCHAR(100)                                    NOT NULL,
    description TEXT                                            NULL,
    is_active   BOOLEAN                  DEFAULT TRUE,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP WITH TIME ZONE                        NULL
);

-- Tabel untuk menyimpan metode pembayaran
CREATE TABLE payment_methods
(
    id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY NOT NULL,
    code       VARCHAR(50) UNIQUE                              NOT NULL, -- 'bank_transfer', 'e_wallet', 'credit_card', 'va'
    name       VARCHAR(100)                                    NOT NULL,
    provider   VARCHAR(100)                                    NOT NULL, -- 'midtrans', 'xendit', 'doku', 'manual'
    currency   varchar(4)                                      NOT NULL DEFAULT 'IDR',
    fee_type   varchar(10)                                              DEFAULT 'none',
    fee_amount DECIMAL(15, 2)                                           DEFAULT 0,
    is_active  BOOLEAN                                                  DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE                                 DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE                                 DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE                        NULL
);

-- Tabel utama untuk transaksi pembayaran
CREATE TABLE payments
(
    id                     BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY NOT NULL,
    transaction_id         VARCHAR(100) UNIQUE                             NOT NULL, -- Generated unique ID
    platform_id            BIGINT REFERENCES platforms (id),
    payment_type_id        BIGINT REFERENCES payment_types (id),
    payment_method_id      BIGINT REFERENCES payment_methods (id),
    amount                 DECIMAL(15, 2)                                  NOT NULL,
    fee_amount             DECIMAL(15, 2)           DEFAULT 0,
    total_amount           DECIMAL(15, 2)                                  NOT NULL, -- amount + fee_amount
    currency               VARCHAR(3)               DEFAULT 'IDR',
    status                 VARCHAR(10)              DEFAULT 'pending',
    customer_id            VARCHAR(100),                                             -- ID customer dari platform terkait
    customer_name          VARCHAR(255),
    customer_email         VARCHAR(255),
    customer_phone         VARCHAR(20),
    reference_id           VARCHAR(100),                                             -- Order ID, Token Request ID, dll
    reference_type         VARCHAR(50),                                              -- 'order', 'token_request', 'wallet_request'
    gateway_transaction_id VARCHAR(255),                                             -- ID dari payment gateway
    gateway_reference      VARCHAR(255),
    gateway_response       JSONB,                                                    -- JSON response dari gateway
    callback_url           VARCHAR(500),
    return_url             VARCHAR(500),
    expired_at             TIMESTAMP,
    paid_at                TIMESTAMP,
    created_at             TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at             TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at             TIMESTAMP WITH TIME ZONE                        NULL
);

-- Create indexes separately for payments table
CREATE INDEX idx_payments_transaction_id ON payments (transaction_id);
CREATE INDEX idx_payments_platform_status ON payments (platform_id, status);
CREATE INDEX idx_payments_customer_id ON payments (customer_id);
CREATE INDEX idx_payments_reference ON payments (reference_id, reference_type);
CREATE INDEX idx_payments_gateway_transaction ON payments (gateway_transaction_id);
CREATE INDEX idx_payments_created_at ON payments (created_at);

-- Tabel untuk menyimpan detail item pembayaran (untuk sales order)
CREATE TABLE payment_items
(
    id               BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY NOT NULL,
    payment_id       BIGINT REFERENCES payments (id) ON DELETE CASCADE,
    item_id          VARCHAR(100), -- SKU atau ID produk
    item_name        VARCHAR(255)                                    NOT NULL,
    item_description TEXT                                            NULL,
    quantity         INTEGER                  DEFAULT 1,
    unit_price       DECIMAL(15, 2)                                  NOT NULL,
    total_price      DECIMAL(15, 2)                                  NOT NULL,
    created_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at       TIMESTAMP WITH TIME ZONE                        NULL
);

-- Tabel untuk menyimpan riwayat status pembayaran
CREATE TABLE payment_status_history
(
    id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY NOT NULL,
    payment_id BIGINT REFERENCES payments (id) ON DELETE CASCADE,
    status     VARCHAR(50)                                     NOT NULL,
    notes      TEXT,
    created_by JSONB, -- system, gateway, admin
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE                        NULL
);

CREATE INDEX idx_payment_status_history_payment_id ON payment_status_history (payment_id);
CREATE INDEX idx_payment_status_history_created_at ON payment_status_history (created_at);

-- Tabel untuk menyimpan log callback dari payment gateway
CREATE TABLE payment_callbacks
(
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY NOT NULL,
    payment_id    BIGINT REFERENCES payments (id),
    gateway_name  VARCHAR(100),
    callback_data JSONB, -- JSON data dari gateway
    response_data JSONB, -- Response yang dikirim kembali
    status        VARCHAR(50),
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMP WITH TIME ZONE                        NULL
);

CREATE INDEX idx_payment_callbacks_payment_id ON payment_callbacks (payment_id);
CREATE INDEX idx_payment_callbacks_created_at ON payment_callbacks (created_at);

-- Tabel untuk konfigurasi per platform
CREATE TABLE platform_configurations
(
    id           BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY NOT NULL,
    platform_id  BIGINT REFERENCES platforms (id),
    config_key   VARCHAR(100)                                    NOT NULL,
    config_value TEXT,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMP WITH TIME ZONE                        NULL,
    UNIQUE (platform_id, config_key)
);

-- ============================================================================
-- ADMIN & USER MANAGEMENT TABLES
-- ============================================================================

-- Tabel untuk menyimpan role/peran admin
CREATE TABLE admin_roles
(
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY NOT NULL,
    code        VARCHAR(50) UNIQUE                              NOT NULL, -- 'super_admin', 'admin', 'operator', 'viewer'
    name        VARCHAR(100)                                    NOT NULL,
    description TEXT,
    permissions JSON,                                                     -- Daftar permissions dalam format JSON
    is_active   BOOLEAN                  DEFAULT TRUE,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP WITH TIME ZONE                        NULL
);

-- Tabel untuk menyimpan admin users
CREATE TABLE admin_users
(
    id                    BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY NOT NULL,
    username              VARCHAR(100) UNIQUE                             NOT NULL,
    email                 VARCHAR(255) UNIQUE                             NOT NULL,
    password_hash         VARCHAR(255)                                    NOT NULL, -- Hashed password
    full_name             VARCHAR(255)                                    NOT NULL,
    phone                 VARCHAR(20),
    avatar_url            VARCHAR(500),
    role_id               BIGINT REFERENCES admin_roles (id),
    is_active             BOOLEAN                  DEFAULT TRUE,
    is_verified           BOOLEAN                  DEFAULT FALSE,
    last_login_at         TIMESTAMP,
    last_login_ip         VARCHAR(45),
    failed_login_attempts INT                      DEFAULT 0,
    locked_until          TIMESTAMP,
    password_changed_at   TIMESTAMP                DEFAULT CURRENT_TIMESTAMP,
    two_factor_enabled    BOOLEAN                  DEFAULT FALSE,
    two_factor_secret     VARCHAR(255),
    created_at            TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at            TIMESTAMP WITH TIME ZONE                        NULL
);

CREATE INDEX idx_admin_users_username ON admin_users (username);
CREATE INDEX idx_admin_users_email ON admin_users (email);
CREATE INDEX idx_admin_users_role_active ON admin_users (role_id, is_active);
CREATE INDEX idx_admin_users_last_login ON admin_users (last_login_at);

-- Tabel untuk menyimpan session admin
CREATE TABLE admin_sessions
(
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY NOT NULL,
    admin_user_id BIGINT REFERENCES admin_users (id) ON DELETE CASCADE,
    session_token VARCHAR(255) UNIQUE                             NOT NULL,
    refresh_token VARCHAR(255) UNIQUE,
    ip_address    VARCHAR(45),
    user_agent    TEXT,
    expires_at    TIMESTAMP                                       NOT NULL,
    is_active     BOOLEAN                     DEFAULT TRUE,
    last_used_at  TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at    TIMESTAMP WITH TIME ZONE    DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP WITH TIME ZONE    DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMP WITH TIME ZONE                        NULL
);

CREATE INDEX idx_admin_sessions_session_token ON admin_sessions (session_token);
CREATE INDEX idx_admin_sessions_refresh_token ON admin_sessions (refresh_token);
CREATE INDEX idx_admin_sessions_admin_user ON admin_sessions (admin_user_id);
CREATE INDEX idx_admin_sessions_expires_at ON admin_sessions (expires_at);

-- Tabel untuk menyimpan log aktivitas admin
CREATE TABLE admin_activity_logs
(
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY NOT NULL,
    admin_user_id BIGINT REFERENCES admin_users (id),
    action        VARCHAR(100)                                    NOT NULL, -- 'login', 'logout', 'create_payment', 'cancel_payment', etc
    resource_type VARCHAR(50),                                              -- 'payment', 'user', 'platform', etc
    resource_id   VARCHAR(100),                                             -- ID dari resource yang diakses
    description   TEXT,
    ip_address    VARCHAR(45),
    user_agent    TEXT,
    request_data  JSON,                                                     -- Data request jika ada
    response_data JSON,                                                     -- Data response jika ada
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMP WITH TIME ZONE                        NULL
);

CREATE INDEX idx_admin_activity_logs_admin_user_action ON admin_activity_logs (admin_user_id, action);
CREATE INDEX idx_admin_activity_logs_resource ON admin_activity_logs (resource_type, resource_id);
CREATE INDEX idx_admin_activity_logs_created_at ON admin_activity_logs (created_at);

-- Tabel untuk menyimpan permission yang bisa diassign ke role
-- CREATE TABLE permissions
-- (
--     id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
--     code        VARCHAR(100) UNIQUE NOT NULL, -- 'payment.create', 'payment.view', 'user.manage', etc
--     name        VARCHAR(255)        NOT NULL,
--     description TEXT,
--     module      VARCHAR(50)         NOT NULL, -- 'payment', 'user', 'platform', 'report'
--     created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
-- );

-- Tabel untuk mapping role dan permission (many-to-many)
-- CREATE TABLE role_permissions
-- (
--     id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
--     role_id       BIGINT REFERENCES admin_roles (id) ON DELETE CASCADE,
--     permission_id BIGINT REFERENCES permissions (id) ON DELETE CASCADE,
--     created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     UNIQUE (role_id, permission_id)
-- );