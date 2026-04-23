package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type Career_Stats struct {
	PlayerID                      int
	SeasonID                      string
	LeagueID                      string
	TeamID                        int
	TeamAbbreviation              string
	PlayerAge                     float64
	GamesPlayed                   int
	GamesStarted                  int
	Minutes                       int
	FieldGoalMade                 int
	FieldGoalAttempts             int
	FieldGoalPercentage           float64
	ThreePointFieldGoalMade       int
	ThreePointFieldGoalAttempts   int
	ThreePointFiledGoalPercentage float64
	FreeThrowMade                 int
	FreeThrowAttempts             int
	FreeThrowPercentage           float64
	OffensiveRebounds             int
	DefensiveRebounds             int
	Rebounds                      int
	Assists                       int
	Steals                        int
	Blocks                        int
	Turnovers                     int
	PersonalFouls                 int
	Points                        int
	PlayerName                    string
}

func main() {
	// 1. Open connection
	connStr := "postgresql://nba:nba@172.21.192.1:5432/nbadb?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 2. Execute SELECT query
	rows, err := db.Query("SELECT * FROM career_stats")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// 3. Iterate over results
	var career_stats []Career_Stats
	for rows.Next() {
		var cS Career_Stats
		// Scan columns into struct fields
		if err := rows.Scan(&cS.PlayerID, &cS.SeasonID, &cS.LeagueID, &cS.TeamID, &cS.TeamAbbreviation, &cS.PlayerAge, &cS.GamesPlayed, &cS.GamesStarted, &cS.Minutes, &cS.FieldGoalMade, &cS.FieldGoalAttempts, &cS.FieldGoalPercentage, &cS.ThreePointFieldGoalMade, &cS.ThreePointFieldGoalAttempts, &cS.ThreePointFiledGoalPercentage, &cS.FreeThrowMade, &cS.FreeThrowAttempts, &cS.FreeThrowPercentage, &cS.OffensiveRebounds, &cS.DefensiveRebounds, &cS.Rebounds, &cS.Assists, &cS.Steals, &cS.Blocks, &cS.Turnovers, &cS.PersonalFouls, &cS.Points, &cS.PlayerName); err != nil {
			log.Fatal(err)
		}
		career_stats = append(career_stats, cS)
	}

	// 4. Check for errors during iteration
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Retrieved %d career stats records\n", len(career_stats))
}
