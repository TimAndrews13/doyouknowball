package game

import (
	"context"
	"fmt"
	"time"

	sqlc "github.com/timandrews/doyouknowball/internal/db/sqlc"
)

const MaxGuesses = 3

type TeamStop struct {
	TeamAbbreviation string  `json:"team_abbreviation"`
	StartYear        string  `json:"start_year"`
	EndYear          string  `json:"end_year"`
	GP               int64   `json:"gp"`
	PPG              float64 `json:"ppg"`
	RPG              float64 `json:"rpg"`
	APG              float64 `json:"apg"`
	SPG              float64 `json:"spg"`
	BPG              float64 `json:"bpg"`
	FGPct            float64 `json:"fg_pct"`
	FG3Pct           float64 `json:"fg3_pct"`
	FTPct            float64 `json:"ft_pct"`
}

type CareerPath struct {
	TeamStops []TeamStop `json:"team_stops"`
}

func BuildCareerPath(stats []sqlc.GetPlayerCareerPathRow) CareerPath {
	// Group stats by team
	type teamData struct {
		games     int64
		points    int64
		rebounds  int64
		assists   int64
		steals    int64
		blocks    int64
		fgm       int64
		fga       int64
		fg3m      int64
		fg3a      int64
		ftm       int64
		fta       int64
		startYear string
		endYear   string
	}

	// Use a slice to preserve order
	teamOrder := []string{}
	teamMap := map[string]*teamData{}

	for _, s := range stats {
		abbr := s.TEAMABBREVIATION.String
		year := s.SEASONID[:4] // e.g. "2020" from "2020-21"

		if _, exists := teamMap[abbr]; !exists {
			teamMap[abbr] = &teamData{startYear: year}
			teamOrder = append(teamOrder, abbr)
		}

		td := teamMap[abbr]
		td.endYear = fmt.Sprintf("%d", mustInt(year)+1)
		td.games += s.GP.Int64
		td.points += s.PTS.Int64
		td.rebounds += s.REB.Int64
		td.assists += s.AST.Int64
		td.steals += s.STL.Int64
		td.blocks += s.BLK.Int64
		td.fgm += s.FGM.Int64
		td.fga += s.FGA.Int64
		td.fg3m += s.FG3M.Int64
		td.fg3a += s.FG3A.Int64
		td.ftm += s.FTM.Int64
		td.fta += s.FTA.Int64
	}

	var teamStops []TeamStop
	for _, abbr := range teamOrder {
		td := teamMap[abbr]
		g := float64(td.games)

		stop := TeamStop{
			TeamAbbreviation: abbr,
			StartYear:        td.startYear,
			EndYear:          td.endYear,
			GP:               td.games,
			PPG:              round(float64(td.points) / g),
			RPG:              round(float64(td.rebounds) / g),
			APG:              round(float64(td.assists) / g),
			SPG:              round(float64(td.steals) / g),
			BPG:              round(float64(td.blocks) / g),
			FGPct:            round(float64(td.fgm) / float64(max64(td.fga, 1)) * 100),
			FG3Pct:           round(float64(td.fg3m) / float64(max64(td.fg3a, 1)) * 100),
			FTPct:            round(float64(td.ftm) / float64(max64(td.fta, 1)) * 100),
		}
		teamStops = append(teamStops, stop)
	}

	return CareerPath{TeamStops: teamStops}
}

func round(f float64) float64 {
	return float64(int(f*10+0.5)) / 10
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func mustInt(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}

func IsGameComplete(guessCount int, correct bool) bool {
	return correct || guessCount >= MaxGuesses
}

func CheckGuess(guess string, playerName string) bool {
	return guess == playerName
}

// Placeholder for future scheduled job
func ScheduleDailyGame(ctx context.Context, queries *sqlc.Queries) error {
	player, err := queries.GetRandomPlayer(ctx)
	if err != nil {
		return fmt.Errorf("failed to get random player: %w", err)
	}

	_, err = queries.CreateDailyGame(ctx, int32(player.PLAYERID))
	if err != nil {
		return fmt.Errorf("failed to create daily game: %w", err)
	}

	return nil
}

func StartDailyGameScheduler(ctx context.Context, queries *sqlc.Queries) {
	go func() {
		for {
			now := time.Now()

			// Calculate time until next midnight EST
			est, _ := time.LoadLocation("America/New_York")
			nowEST := now.In(est)
			nextMidnight := time.Date(nowEST.Year(), nowEST.Month(), nowEST.Day()+1, 0, 0, 0, 0, est)
			duration := nextMidnight.Sub(now)

			select {
			case <-time.After(duration):
				if err := ScheduleDailyGame(ctx, queries); err != nil {
					fmt.Println("scheduler error:", err)
				} else {
					fmt.Println("daily game scheduled successfully")
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	fmt.Println("daily game scheduler started")
}
