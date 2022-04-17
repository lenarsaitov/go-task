DROP TABLE IF EXISTS cards_history;
DROP TABLE IF EXISTS cards;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
       user_id          SERIAL PRIMARY KEY,
       user_full_name   varchar(100) NOT NULL,
       create_time      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE cards (
       card_id       SERIAL PRIMARY KEY,
       balance       BIGINT NOT NULL DEFAULT 0,
       create_time   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
       user_id       INT REFERENCES users (user_id)
);