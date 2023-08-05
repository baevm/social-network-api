CREATE TABLE IF NOT EXISTS followers (
    user_id bigserial REFERENCES users (id),
    follower_id bigserial REFERENCES users (id),
    followed_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, follower_id),
    UNIQUE (user_id, follower_id)
);