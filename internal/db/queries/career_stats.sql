-- name: GetPlayerCareerPath :many
SELECT 
    "TEAM_ABBREVIATION",
    "SEASON_ID",
    "PLAYER_AGE",
    "GP",
    "MIN",
    "FGM",
    "FGA",
    "FG_PCT",
    "FG3M",
    "FG3A",
    "FG3_PCT",
    "FTM",
    "FTA",
    "FT_PCT",
    "OREB",
    "DREB",
    "REB",
    "AST",
    "STL",
    "BLK",
    "TOV",
    "PF",
    "PTS"
FROM career_stats
WHERE "PLAYER_ID" = $1
ORDER BY "SEASON_ID" ASC;

-- name: GetAllPlayerNames :many
SELECT DISTINCT "PLAYER_ID", "PLAYER_NAME"
FROM career_stats
ORDER BY "PLAYER_NAME" ASC;

-- name: GetRandomPlayer :one
SELECT DISTINCT ON (random()) "PLAYER_ID", "PLAYER_NAME"
FROM career_stats
ORDER BY random()
LIMIT 1;

-- name: GetPlayerByID :one
SELECT DISTINCT "PLAYER_ID", "PLAYER_NAME"
FROM career_stats
WHERE "PLAYER_ID" = $1;

-- name: SearchPlayers :many
SELECT DISTINCT "PLAYER_ID", "PLAYER_NAME"
FROM career_stats
WHERE LOWER("PLAYER_NAME") LIKE LOWER($1)
ORDER BY "PLAYER_NAME" ASC
LIMIT 10;