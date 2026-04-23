-- name: CreateDailyGame :one
INSERT INTO daily_game (player_id, game_date)
VALUES ($1, $2)
RETURNING *;

-- name: GetTodayGame :one
SELECT * FROM daily_game
WHERE game_date = CURRENT_DATE;

-- name: GetDailyGameByID :one
SELECT * FROM daily_game
WHERE id = $1;