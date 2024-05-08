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
    role VARCHAR(10) NOT NULL, 
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE products (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    description TEXT,
    made_in VARCHAR(100) NOT NULL,
    count INT NOT NULL,
    cost FLOAT NOT NULL,
    discount FLOAT NOT NULL DEFAULT 0,
    color VARCHAR(20),
    size INT NOT NULL,
    age_min INT NOT NULL DEFAULT 1,
    age_max INT,
    for_gender VARCHAR(6) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE media (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL, 
    image_url TEXT NOT NULL,
    file_name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (product_id) REFERENCES products (id)
);

CREATE TABLE orders (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    product_id UUID NOT NULL,
    status VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (product_id) REFERENCES products (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE stars (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL,
    user_id UUID NOT NULL,
    star INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    FOREIGN KEY (product_id) REFERENCES products (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE wishlist (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    FOREIGN KEY (product_id) REFERENCES products (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE saves (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    FOREIGN KEY (product_id) REFERENCES products (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE comments (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL,
    user_id UUID NOT NULL,
    comment TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    FOREIGN KEY (product_id) REFERENCES products (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

INSERT INTO users (id, first_name, last_name, email, phone_number, password, gender, age, refresh, role) 
VALUES 
('19d16003-586a-4190-92ee-ab0c45504023', 'Xasan', 'Nosirov', 'xasannosirov094@gmail.com', '+998944970514', '$2a$10$VOukMtTpUxICddVOCTJJou594V0cZ4zbRVN9smlrcMrH6i4AjqrbK', 'male', 18, NULL, 'admin'),
('19d16003-586a-4190-92ee-ab0c45504024', 'Alisher', 'Botirov', 'xasannosirov094@gmail.com', '+998998887766', '$2a$10$VOukMtTpUxICddVOCTJJou594V0cZ4zbRVN9smlrcMrH6i4AjqrbK', 'male', 23, NULL, 'worker'),
('19d16003-586a-4190-92ee-ab0c45504025', 'Jasur', 'Abudllaev', 'xasannosirov094@gmail.com', '+998997778866', '$2a$10$VOukMtTpUxICddVOCTJJou594V0cZ4zbRVN9smlrcMrH6i4AjqrbK', 'male', 25, NULL, 'user');
