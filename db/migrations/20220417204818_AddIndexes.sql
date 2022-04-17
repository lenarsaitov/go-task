
-- +goose Up
CREATE INDEX idx_users_id ON users(user_id);
CREATE INDEX idx_cards_id ON cards(card_id);

-- +goose Down
DROP INDEX IF EXISTS idx_users_id;
DROP INDEX IF EXISTS idx_cards_id;
