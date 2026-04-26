-- +goose Up
CREATE TABLE daily_game (
    id          SERIAL PRIMARY KEY,
    player_id   INT NOT NULL,
    game_date   TIMESTAMP UNIQUE NOT NULL DEFAULT NOW(),
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE daily_game;