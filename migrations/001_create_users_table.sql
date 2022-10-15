CREATE TABLE IF NOT EXISTS users
(
    id       SERIAL PRIMARY KEY,
    login    varchar(255) NOT NULL UNIQUE,
    password varchar(255) NOT NULL,
    created_at timestamp,
    updated_at timestamp
);
