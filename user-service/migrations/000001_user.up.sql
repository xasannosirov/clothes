CREATE TABLE users (
    id UUID PRIMARY KEY,
    first_name VARCHAR(64) NOT NULL,
    last_name VARCHAR(64) NOT NULL,
    email TEXT NOT NULL,
    phone_number VARCHAR(20),
    password TEXT NOT NULL,
    gender VARCHAR(6) NOT NULL,
    age INT,
    refresh TEXT,
    role VARCHAR(10) NOT NULL, -- admin, worker, user
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);