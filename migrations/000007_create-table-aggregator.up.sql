CREATE TABLE IF NOT EXISTS aggregators
(
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY NOT NULL,
    name        VARCHAR(50) UNIQUE                              NOT NULL,
    slug        varchar(50) UNIQUE                              NOT NULL,
    description TEXT,
    is_active   BOOLEAN                  DEFAULT TRUE,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP WITH TIME ZONE                        NULL
)