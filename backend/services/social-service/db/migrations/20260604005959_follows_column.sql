-- +goose Up

CREATE TABLE follows
(
    follower_id UUID NOT NULL,
    followee_id UUID NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (follower_id, followee_id),

    CONSTRAINT check_not_self_follow
        CHECK (follower_id <> followee_id)
);

CREATE INDEX idx_follows_follower_id ON follows (follower_id);
CREATE INDEX idx_follows_followee_id ON follows (followee_id);

-- +goose Down

DROP TABLE IF EXISTS follows;