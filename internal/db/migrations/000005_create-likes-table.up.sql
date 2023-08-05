CREATE TABLE IF NOT EXISTS post_like (
    post_id bigserial NOT NULL REFERENCES posts (id) ON DELETE CASCADE,
    user_id bigserial NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    PRIMARY KEY (post_id, user_id)
);

CREATE INDEX ON post_like (post_id);