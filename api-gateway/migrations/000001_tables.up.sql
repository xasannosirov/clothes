CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    first_name VARCHAR(64) NOT NULL,
    last_name VARCHAR(64) NOT NULL,
    email TEXT NOT NULL,
    phone_number VARCHAR(20),
    password TEXT,
    gender VARCHAR(6),
    age INT,
    refresh TEXT,
    role VARCHAR(10) NOT NULL, 
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS category (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    category_id UUID NOT NULL,
    description TEXT,
    made_in VARCHAR(100) NOT NULL,
    count INT NOT NULL,
    cost FLOAT NOT NULL,
    discount FLOAT NOT NULL DEFAULT 0,
    color VARCHAR(20) [],
    size VARCHAR(5) [] NOT NULL,
    age_min INT NOT NULL DEFAULT 1,
    age_max INT,
    for_gender VARCHAR(6) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (category_id) REFERENCES category (id)
);

CREATE TABLE IF NOT EXISTS media (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL, 
    image_url TEXT NOT NULL,
    file_name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (product_id) REFERENCES products (id)
);

CREATE TABLE IF NOT EXISTS wishlist (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    FOREIGN KEY (product_id) REFERENCES products (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- Create comments table
CREATE TABLE comments (
    id UUID PRIMARY KEY, 
    owner_id UUID NOT NULL,
    product_id UUID NOT NULL, 
    message TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (owner_id) REFERENCES users (id),
    FOREIGN KEY (product_id) REFERENCES products (id)
);

 
CREATE TABLE baskets (
    user_id UUID PRIMARY KEY, 
    product_id UUID[],
    FOREIGN KEY (user_id) REFERENCES users (id)
);

INSERT INTO users (id, first_name, last_name, email, phone_number, password, gender, age, refresh, role) 
VALUES 
('19d16003-586a-4190-92ee-ab0c45504023', 'Xasan', 'Nosirov', 'xasannosirov094@gmail.com', '+998944970514', '$2a$10$VOukMtTpUxICddVOCTJJou594V0cZ4zbRVN9smlrcMrH6i4AjqrbK', 'male', 18, NULL, 'admin');
