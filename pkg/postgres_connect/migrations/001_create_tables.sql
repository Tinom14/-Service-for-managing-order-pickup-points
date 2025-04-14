-- +migrate Up
CREATE TABLE users
(
    id            SERIAL PRIMARY KEY,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255)        NOT NULL,
    role          VARCHAR(50)         NOT NULL CHECK (role IN ('employee', 'moderator'))
);

CREATE TABLE pvz
(
    id                SERIAL PRIMARY KEY,
    city              VARCHAR(100) NOT NULL CHECK (city IN ('Москва', 'Санкт-Петербург', 'Казань')),
    registration_date TIMESTAMP DEFAULT NOW()
);

CREATE TABLE receptions
(
    id         SERIAL PRIMARY KEY,
    pvz_id     INT REFERENCES pvz (id),
    status     VARCHAR(20)             NOT NULL CHECK (status IN ('in_progress', 'closed')),
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    FOREIGN KEY (pvz_id) REFERENCES pvz (id) ON DELETE CASCADE
);

CREATE TABLE products
(
    id       SERIAL PRIMARY KEY,
    type     VARCHAR(50)             NOT NULL CHECK (type IN ('электроника', 'одежда', 'обувь')),
    added_at TIMESTAMP DEFAULT NOW() NOT NULL
);

CREATE TABLE reception_products
(
    reception_id INT NOT NULL,
    product_id   INT NOT NULL,
    PRIMARY KEY (reception_id, product_id),
    FOREIGN KEY (reception_id) REFERENCES receptions (id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS receptions;
DROP TABLE IF EXISTS pvz;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS reception_products;