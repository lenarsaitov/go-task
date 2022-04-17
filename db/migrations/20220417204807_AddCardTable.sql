
-- +goose Up
CREATE TABLE cards (
                       card_id       SERIAL PRIMARY KEY,
                       balance       BIGINT NOT NULL DEFAULT 0,
                       create_time   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
                       user_id       INT REFERENCES users (user_id)
);

-- +goose Down
DROP TABLE IF EXISTS cards;

