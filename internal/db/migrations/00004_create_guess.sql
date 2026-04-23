-- +goose Up
CREATE TABLE guesses (
    id            SERIAL PRIMARY KEY,
    user_id       INT NOT NULL REFERENCES users(id),
    daily_game_id INT NOT NULL REFERENCES daily_game(id),
    guess         VARCHAR(255) NOT NULL,
    is_correct    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE guesses;