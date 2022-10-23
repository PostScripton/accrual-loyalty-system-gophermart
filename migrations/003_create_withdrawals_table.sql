CREATE TABLE withdrawals
(
    id         SERIAL PRIMARY KEY,
    user_id    INT          NOT NULL REFERENCES users (id),
    number     VARCHAR(255) NOT NULL UNIQUE,
    sum        DECIMAL      NULL,
    created_at TIMESTAMP DEFAULT now()
);
