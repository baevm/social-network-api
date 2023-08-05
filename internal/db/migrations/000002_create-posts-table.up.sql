CREATE TABLE IF NOT EXISTS posts (
    id bigserial PRIMARY KEY,
    user_id bigserial NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    body text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS post_images (
    id bigserial PRIMARY KEY,
    post_id bigserial NOT NULL REFERENCES posts (id),
    url text NOT NULL
);

CREATE INDEX ON post_images (post_id);