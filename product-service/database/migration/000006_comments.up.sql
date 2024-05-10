CREATE TABLE IF NOT EXISTS COMMENTS(
    ID UUID PRIMARY KEY,
    PRODUCT_ID UUID NOT NULL,
    USER_ID UUID NOT NULL,
    COMMENT TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    FOREIGN KEY (product_id) REFERENCES products (id)
);