CREATE TABLE IF NOT EXISTS product
(
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    price NUMERIC(6, 2) NOT NULL,
    quantity INT NOT NULL
    );
