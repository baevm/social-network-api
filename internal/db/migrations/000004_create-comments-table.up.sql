CREATE TABLE IF NOT EXISTS comments (
    id bigserial PRIMARY KEY,
    post_id bigserial NOT NULL REFERENCES posts (id) ON DELETE CASCADE,
    user_id bigserial NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    body text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS comment_images (
    id bigserial PRIMARY KEY,
    comment_id bigserial NOT NULL REFERENCES comments (id) ON DELETE CASCADE,
    url text NOT NULL
);

CREATE INDEX ON comment_images (comment_id);