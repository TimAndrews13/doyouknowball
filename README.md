# DoYouKnowBall 🏀

A daily NBA player guessing game inspired by Wordle. Each day a random NBA player is selected and you must identify them from their career path — the teams they played for, the years they were there, and their stats along the way.

---

## How to Play

1. Log in or create an account
2. A player's career path is shown — the teams they played for, the seasons, and their per-game stats for each stint
3. Search for a player and submit your guess
4. You have **3 guesses** to identify the player
5. A new player is available every day at 5am EST

---

## Tech Stack

**Backend**
- [Go](https://golang.org/) — HTTP server and API
- [PostgreSQL](https://www.postgresql.org/) — Database
- [Docker](https://www.docker.com/) — Database containerization
- [Goose](https://github.com/pressly/goose) — Database migrations
- [sqlc](https://sqlc.dev/) — Type-safe SQL query generation
- [golang-jwt](https://github.com/golang-jwt/jwt) — JWT authentication
- [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) — Password hashing
- [godotenv](https://github.com/joho/godotenv) — Environment variable management

**Frontend**
- Go `html/template` — Server-side HTML rendering
- Vanilla JavaScript (embedded in templates) — API calls, autocomplete dropdown, game interactions, dark mode toggle, logout
- CSS custom properties — Light/dark mode theming

**Data**
- [nba_api](https://github.com/swar/nba_api) — Python library for pulling NBA stats
- [SQLAlchemy](https://www.sqlalchemy.org/) — Python database ORM
- [pandas](https://pandas.pydata.org/) — Data manipulation

---

## Project Structure

```
doyouknowball/
├── cmd/
│   └── server/
│       └── main.go           # Entry point, routes, handlers
├── internal/
│   ├── auth/
│   │   └── auth.go           # JWT generation, password hashing
│   ├── game/
│   │   └── game.go           # Career path logic, guess checking, daily scheduler
│   └── db/
│       ├── db.go             # Database connection
│       ├── migrations/       # Goose migration files
│       ├── queries/          # sqlc SQL query files
│       └── sqlc/             # sqlc generated Go code
├── data/
│   ├── pull_player_stats.py  # Fetches NBA player stats into DB
│   └── download_logos.py     # Downloads NBA team logos (only needed if logos are missing)
├── web/
│   ├── templates/
│   │   ├── base.html         # Shared layout, nav, dark mode toggle, logout
│   │   ├── login.html        # Login/register page
│   │   └── game.html         # Main game page
│   └── static/
│       ├── css/
│       │   └── style.css     # Styles, light/dark theme
│       └── images/
│           └── logos/        # NBA team logo PNGs (included in repo)
├── .env                      # Environment variables (not committed)
├── docker-compose.yml        # PostgreSQL container
├── go.mod
└── go.sum
```

---

## Getting Started

### Prerequisites

- [Go 1.21+](https://golang.org/dl/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [Python 3.12+](https://www.python.org/)
- [Goose](https://github.com/pressly/goose) — `go install github.com/pressly/goose/v3/cmd/goose@latest`
- [sqlc](https://sqlc.dev/) — `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`

### Setup

**1. Clone the repo**
```bash
git clone https://github.com/timandrews/doyouknowball.git
cd doyouknowball
```

**2. Create a `.env` file**
```
DB_URL=postgresql://nba:nba@localhost:5432/nbadb?sslmode=disable
JWT_SECRET=your-secret-key
```

**3. Start the database**
```bash
docker compose up -d
```

**4. Run migrations**
```bash
cd internal/db/migrations
goose postgres $DB_URL up
```

**5. Set up Python environment and fetch player data**
```bash
python3 -m venv venv
source venv/bin/activate
pip install nba_api pandas sqlalchemy psycopg2-binary python-dotenv
python3 data/pull_player_stats.py
```

**6. Run the server**
```bash
go run ./cmd/server/
```

**7. Set up the daily game**
```bash
curl -X POST http://localhost:8080/api/game/setup
```

Visit `http://localhost:8080` in your browser.

> **Note:** Team logos are already included in the repo at `web/static/images/logos/`. You only need to run `data/download_logos.py` if logos are missing.

---

## API Endpoints

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/register` | No | Create a new account |
| POST | `/api/login` | No | Log in and receive a JWT |
| GET | `/api/game/today` | Yes | Get today's career path, guess count, and previous guesses |
| POST | `/api/game/guess` | Yes | Submit a guess |
| POST | `/api/game/setup` | No | Set up today's daily game (manual trigger) |
| GET | `/api/players/search?q=` | Yes | Search for players by name |

---

## Data Pipeline

Player stats are fetched from the NBA API using Python and stored in PostgreSQL. To refresh player data:

```bash
source venv/bin/activate
python3 data/pull_player_stats.py
```

The script deletes and re-inserts stats for each player so trades and new seasons are always up to date.

---

## Daily Game Scheduler

A new player is automatically selected every day at **5am EST**. The scheduler runs as a background goroutine when the server starts. For local development you can manually trigger a new game at any time:

```bash
curl -X POST http://localhost:8080/api/game/setup
```

---

## Roadmap

- [x] Scheduled daily game picker (5am EST)
- [x] Persist guess history across page refreshes
- [x] Handle already-completed games on page load
- [x] Logout button
- [ ] Mobile optimization
- [ ] Player headshot reveal after correct guess
- [ ] User stats dashboard (daily streak, win %, games played, guess distribution)
- [ ] Friend groups with leaderboards and group history tracking
- [ ] Production deployment