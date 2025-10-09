-- Active: 1744740356603@@127.0.0.1@5432@artworks

-- SELECT TABLE_NAME
-- FROM INFORMATION_SCHEMA.TABLES

CREATE TABLE users (
    id UUID PRIMARY KEY,
    login VARCHAR(50) NOT NULL UNIQUE,
    hashed_password VARCHAR(255) NOT NULL
);
ALTER TABLE users ADD CONSTRAINT user_empty_check 
    CHECK(login != '' AND hashed_password != ''); 

CREATE TABLE shops (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(500) NOT NULL,
    user_id UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    update_time TIMESTAMP NOT NULL
);

ALTER TABLE shops ADD CONSTRAINT shop_empty_check 
    CHECK(title != ''); 

CREATE TABLE posts (
    id UUID PRIMARY KEY,
    description VARCHAR(500) NOT NULL,
    publication_time TIMESTAMP NOT NULL,
    shop_id UUID NOT NULL,
    FOREIGN KEY (shop_id) REFERENCES shops(id)
);

CREATE TABLE products (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(500) NOT NULL,
    cost INT,
    shop_id UUID NOT NULL,
    FOREIGN KEY (shop_id) REFERENCES shops(id),
    update_time TIMESTAMP NOT NULL
);

ALTER TABLE products ADD CONSTRAINT product_empty_check 
    CHECK(title != ''); 

CREATE TABLE categories (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(500) NOT NULL
);

ALTER TABLE categories ADD CONSTRAINT category_empty_check 
    CHECK(title != ''); 

CREATE TABLE product_category (
    product_id UUID NOT NULL,
    category_id UUID NOT NULL,
    PRIMARY KEY (product_id, category_id),
    CONSTRAINT product_id_notnull CHECK (product_id IS NOT NULL),
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    CONSTRAINT category_id_notnull CHECK (category_id IS NOT NULL),
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

