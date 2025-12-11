CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL PRIMARY KEY,
    email         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    full_name     TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS devices
(
    id            SERIAL PRIMARY KEY,
    device_id     TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at    TIMESTAMP DEFAULT now()
);
