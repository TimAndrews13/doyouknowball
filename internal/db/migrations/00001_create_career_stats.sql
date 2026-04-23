-- +goose Up
CREATE TABLE career_stats (
    "PLAYER_ID"         BIGINT          NOT NULL,
    "SEASON_ID"         TEXT            NOT NULL,
    "LEAGUE_ID"         TEXT,
    "TEAM_ID"           BIGINT          NOT NULL,
    "TEAM_ABBREVIATION" TEXT,
    "PLAYER_AGE"        DOUBLE PRECISION,
    "GP"                BIGINT,
    "GS"                BIGINT,
    "MIN"               BIGINT,
    "FGM"               BIGINT,
    "FGA"               BIGINT,
    "FG_PCT"            DOUBLE PRECISION,
    "FG3M"              BIGINT,
    "FG3A"              BIGINT,
    "FG3_PCT"           DOUBLE PRECISION,
    "FTM"               BIGINT,
    "FTA"               BIGINT,
    "FT_PCT"            DOUBLE PRECISION,
    "OREB"              BIGINT,
    "DREB"              BIGINT,
    "REB"               BIGINT,
    "AST"               BIGINT,
    "STL"               BIGINT,
    "BLK"               BIGINT,
    "TOV"               BIGINT,
    "PF"                BIGINT,
    "PTS"               BIGINT,
    "PLAYER_NAME"       TEXT,

    PRIMARY KEY ("PLAYER_ID", "SEASON_ID", "TEAM_ID")
);

-- +goose Down
DROP TABLE career_stats;