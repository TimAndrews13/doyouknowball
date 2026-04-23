-- name: CreateGuess :one
INSERT INTO guesses (user_id, daily_game_id, guess, is_correct)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetGuessesByUserAndGame :many
SELECT * FROM guesses
WHERE user_id = $1
AND daily_game_id = $2
ORDER BY created_at ASC;

-- name: GetGuessCount :one
SELECT COUNT(*) FROM guesses
WHERE user_id = $1
AND daily_game_id = $2;