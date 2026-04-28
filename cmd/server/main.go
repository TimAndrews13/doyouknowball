package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/timandrews/doyouknowball/internal/auth"
	db "github.com/timandrews/doyouknowball/internal/db"
	sqlc "github.com/timandrews/doyouknowball/internal/db/sqlc"
	"github.com/timandrews/doyouknowball/internal/game"

	_ "github.com/lib/pq"
)

type App struct {
	db      *sql.DB
	queries *sqlc.Queries
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}

	conn := db.Connect(os.Getenv("DB_URL"))
	defer conn.Close()

	app := &App{
		db:      conn,
		queries: sqlc.New(conn),
	}

	ctx := context.Background()
	game.StartDailyGameScheduler(ctx, app.queries)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	http.HandleFunc("/", app.handleHome)
	http.HandleFunc("/game", app.handleGamePage)
	http.HandleFunc("/api/login", app.handleLogin)
	http.HandleFunc("/api/register", app.handleRegister)
	http.HandleFunc("/api/game/today", app.handleTodayGame)
	http.HandleFunc("/api/game/guess", app.handleGuess)
	http.HandleFunc("/api/game/setup", app.handleSetupDailyGame)
	http.HandleFunc("/api/players/search", app.handlePlayerSearch)
	http.HandleFunc("/api/stats", app.handleStats)

	fmt.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (a *App) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	hashed, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	user, err := a.queries.CreateUser(r.Context(), sqlc.CreateUserParams{
		Username: req.Username,
		Email:    req.Email,
		Password: hashed,
	})
	if err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	token, err := auth.GenerateToken(int(user.ID), user.Username)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (a *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := a.queries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if !auth.CheckPassword(req.Password, user.Password) {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(int(user.ID), user.Username)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (a *App) handleSetupDailyGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := game.ScheduleDailyGame(r.Context(), a.queries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "daily game set up successfully"})
}

func (a *App) handleTodayGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Validate JWT
	claims, err := validateRequest(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	dailyGame, err := a.queries.GetTodayGame(r.Context())
	if err != nil {
		http.Error(w, "no game found for today", http.StatusNotFound)
		return
	}

	stats, err := a.queries.GetPlayerCareerPath(r.Context(), int64(dailyGame.PlayerID))
	if err != nil {
		http.Error(w, "failed to get career path", http.StatusInternalServerError)
		return
	}

	careerPath := game.BuildCareerPath(stats)

	// Get user's guess count
	guessCount, err := a.queries.GetGuessCount(r.Context(), sqlc.GetGuessCountParams{
		UserID:      int32(claims.UserID),
		DailyGameID: dailyGame.ID,
	})
	if err != nil {
		http.Error(w, "failed to get guess count", http.StatusInternalServerError)
		return
	}

	guesses, err := a.queries.GetGuessesByUserAndGame(r.Context(), sqlc.GetGuessesByUserAndGameParams{
		UserID:      int32(claims.UserID),
		DailyGameID: dailyGame.ID,
	})
	if err != nil {
		http.Error(w, "failed to get guesses", http.StatusInternalServerError)
		return
	}

	type guessResult struct {
		Guess     string `json:"guess"`
		IsCorrect bool   `json:"is_correct"`
	}

	previousGuesses := make([]guessResult, len(guesses))
	for i, g := range guesses {
		previousGuesses[i] = guessResult{
			Guess:     g.Guess,
			IsCorrect: g.IsCorrect,
		}
	}

	//Check if game is already complete
	isComplete := false
	answer := ""
	if int(guessCount) >= game.MaxGuesses {
		isComplete = true
	}
	for _, g := range guesses {
		if g.IsCorrect {
			isComplete = true
		}
	}

	if isComplete {
		player, err := a.queries.GetPlayerByID(r.Context(), int64(dailyGame.PlayerID))
		if err != nil {
			http.Error(w, "failed to get player", http.StatusInternalServerError)
			return
		}
		answer = player.PLAYERNAME.String
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"career_path":       careerPath,
		"guesses_remaining": game.MaxGuesses - int(guessCount),
		"previous_guesses":  previousGuesses,
		"is_complete":       isComplete,
		"answer":            answer,
		"player_id":         dailyGame.PlayerID,
	})
}

func (a *App) handleGuess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, err := validateRequest(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Guess string `json:"guess"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	dailyGame, err := a.queries.GetTodayGame(r.Context())
	if err != nil {
		http.Error(w, "no game found for today", http.StatusNotFound)
		return
	}

	guessCount, err := a.queries.GetGuessCount(r.Context(), sqlc.GetGuessCountParams{
		UserID:      int32(claims.UserID),
		DailyGameID: dailyGame.ID,
	})
	if err != nil {
		http.Error(w, "failed to get guess count", http.StatusInternalServerError)
		return
	}

	if int(guessCount) >= game.MaxGuesses {
		http.Error(w, "no guesses remaining", http.StatusForbidden)
		return
	}

	player, err := a.queries.GetPlayerByID(r.Context(), int64(dailyGame.PlayerID))
	if err != nil {
		http.Error(w, "failed to get player", http.StatusInternalServerError)
		return
	}

	isCorrect := game.CheckGuess(req.Guess, player.PLAYERNAME.String)

	_, err = a.queries.CreateGuess(r.Context(), sqlc.CreateGuessParams{
		UserID:      int32(claims.UserID),
		DailyGameID: dailyGame.ID,
		Guess:       req.Guess,
		IsCorrect:   isCorrect,
	})
	if err != nil {
		http.Error(w, "failed to save guess", http.StatusInternalServerError)
		return
	}

	remaining := game.MaxGuesses - int(guessCount) - 1
	response := map[string]interface{}{
		"correct":           isCorrect,
		"guesses_remaining": remaining,
		"player_id":         dailyGame.PlayerID,
	}

	if isCorrect || remaining == 0 {
		response["answer"] = player.PLAYERNAME.String
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func validateRequest(r *http.Request) (*auth.Claims, error) {
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		return nil, fmt.Errorf("missing token")
	}
	// Strip "Bearer " prefix if present
	if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
		tokenStr = tokenStr[7:]
	}
	return auth.ValidateToken(tokenStr)
}

func (a *App) handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/base.html", "web/templates/login.html"))
	tmpl.ExecuteTemplate(w, "base", nil)
}

func (a *App) handleGamePage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/base.html", "web/templates/game.html"))
	tmpl.ExecuteTemplate(w, "base", nil)
}

func (a *App) handlePlayerSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		json.NewEncoder(w).Encode([]struct{}{})
		return
	}

	players, err := a.queries.SearchPlayers(r.Context(), "%"+q+"%")
	if err != nil {
		http.Error(w, "search failed", http.StatusInternalServerError)
		return
	}

	type playerResult struct {
		PlayerID   int64  `json:"player_id"`
		PlayerName string `json:"player_name"`
	}

	results := make([]playerResult, len(players))
	for i, p := range players {
		results[i] = playerResult{
			PlayerID:   p.PLAYERID,
			PlayerName: p.PLAYERNAME.String,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (a *App) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, err := validateRequest(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	history, err := a.queries.GetUserGameHistory(r.Context(), int32(claims.UserID))
	if err != nil {
		http.Error(w, "failed to get history", http.StatusInternalServerError)
		return
	}

	// Calculate stats
	gamesPlayed := len(history)
	wins := 0
	guessDist := map[int]int{1: 0, 2: 0, 3: 0}
	playStreak := 0
	winStreak := 0
	maxWinStreak := 0

	// History is orderd DESC so we calculate streaks from most recent
	estLoc, _ := time.LoadLocation("America/New_York")
	today := time.Now().In(estLoc).Truncate(24 * time.Hour)

	for i, g := range history {
		wonVal, _ := g.Won.(bool)
		if wonVal {
			wins++
			guessesUsed := int(g.GuessesUsed)
			if guessesUsed >= 1 && guessesUsed <= 3 {
				guessDist[guessesUsed]++
			}
		}

		// Play streak - check if consecutive days
		gameDate := g.GameDate.In(estLoc).Truncate(24 * time.Hour)
		expectedDate := today.AddDate(0, 0, -i)
		if gameDate.Equal(expectedDate) {
			playStreak++
		} else {
			break
		}
	}

	// Win streak - separate pass
	for _, g := range history {
		wonVal, _ := g.Won.(bool)
		if wonVal {
			winStreak++
			if winStreak > maxWinStreak {
				maxWinStreak = winStreak
			}
		} else {
			break
		}
	}

	winPct := 0.0
	if gamesPlayed > 0 {
		winPct = math.Round(float64(wins)/float64(gamesPlayed)*100*10) / 10
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"games_played":   gamesPlayed,
		"win_pct":        winPct,
		"play_streak":    playStreak,
		"win_streak":     winStreak,
		"max_win_streak": maxWinStreak,
		"guess_dist":     guessDist,
	})
}
