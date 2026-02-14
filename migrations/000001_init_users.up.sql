CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    user_uid UUID NOT NULL UNIQUE,
    username VARCHAR(30) UNIQUE,
    avatar_uid UUID
);