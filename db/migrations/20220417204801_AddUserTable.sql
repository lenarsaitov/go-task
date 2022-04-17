
-- +goose Up
CREATE TABLE users (
                       user_id          SERIAL PRIMARY KEY,
                       user_full_name   varchar(100) NOT NULL,
                       create_time      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS users;

