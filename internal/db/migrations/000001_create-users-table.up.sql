CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    email citext UNIQUE NOT NULL,
    password_hash text NOT NULL,
    first_name text NOT NULL,
    last_name text NOT NULL,
    username text UNIQUE NOT NULL,
    avatar text,
    birthdate DATE,
    activated bool NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);