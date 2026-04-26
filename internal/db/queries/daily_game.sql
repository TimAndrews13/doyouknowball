-- name: CreateDailyGame :one
INSERT INTO daily_game (player_id, game_date)
VALUES ($1, CURRENT_DATE)
RETURNING *;

-- name: GetTodayGame :one
SELECT id, player_id, game_date, created_at FROM daily_game
WHERE game_date::date = CURRENT_DATE;

-- name: GetDailyGameByID :one
SELECT * FROM daily_game
WHERE id = $1;