CREATE TABLE IF NOT EXISTS SAVES(
    ID UUID PRIMARY KEY,
    PRODUCT_ID UUID NOT NULL,
    USER_ID UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    FOREIGN KEY (product_id) REFERENCES products (id)
);

INSERT INTO SAVES(id, product_id, USER_ID) VALUES(
    '124e6c9c-4e68-4833-9753-9bf74608ee1d', 'f0b91382-3618-429d-8885-17162af3e5df', '5623ea51-e785-42a6-b2fc-ec5c231043c6'
);