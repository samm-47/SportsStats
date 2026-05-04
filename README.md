## Tech Stack

- **Go** — REST API backend, Elo rating engine, SQLite integration
- **Python** — CLI client, batch ingestion (CSV/JSON), Matplotlib visualizations
- **SQLite** — Persistent storage for match history and team stats

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/match` | Record a match result |
| GET | `/team/{name}` | Get stats for a team |
| GET | `/predict/{teamA}/{teamB}` | Get win probabilities for a matchup |

## Getting Started

### Prerequisites
- Go 1.18+
- Python 3.8+
- gcc (required for go-sqlite3)

### Run the Backend

```bash
cd backend
go mod tidy
go run .
```

Server runs at `http://localhost:8080`

### Run the Python Client

```bash
cd client
pip install requests matplotlib
python client.py
```

## How the Rating System Works

Team ratings are initialized at **1500**. After each match, both teams' ratings update based on the Elo formula:

- A team that beats a higher-rated opponent gains more points
- A team that loses to a lower-rated opponent loses more points
- Draws shift ratings toward the middle

The `/predict` endpoint uses current ratings to return a live win probability for any matchup — not a static lookup, but a real-time calculation.
