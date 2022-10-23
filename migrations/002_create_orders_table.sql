CREATE TABLE orders
(
    id         SERIAL PRIMARY KEY,
    user_id    INT          NOT NULL REFERENCES users (id),
    number     VARCHAR(255) NOT NULL UNIQUE,
    status     VARCHAR(255) NOT NULL,
    accrual    DECIMAL      NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);
