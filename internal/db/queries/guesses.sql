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

-- name: GetUserGameHistory :many
SELECT
    dg.id as game_id,
    dg.game_date,
    COALESCE(bool_or(g.is_correct), false) as won,
    COUNT(g.id) as guesses_used
FROM daily_game dg
JOIN guesses g ON g.daily_game_id = dg.id AND g.user_id = $1
GROUP BY dg.id, dg.game_date
ORDER BY dg.game_date DESC;