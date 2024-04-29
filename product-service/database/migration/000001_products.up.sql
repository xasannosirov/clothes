CREATE TABLE IF NOT EXISTS products(
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    description TEXT,
    made_in VARCHAR(100) NOT NULL,
    count INT NOT NULL,
    cost FLOAT NOT NULL,
    discount FLOAT NOT NULL DEFAULT 0,
    color VARCHAR(20),
    size INT,
    age_min INT,
    age_max INT,
    temperature_min INT,
    temperature_max INT,
    for_gender VARCHAR(6),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);