// data.go
package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Match struct {
	TeamA  string `json:"team_a"`
	TeamB  string `json:"team_b"`
	ScoreA int    `json:"score_a"`
	ScoreB int    `json:"score_b"`
}

type TeamStats struct {
	Wins   int
	Losses int
	Draws  int
	Rating float64
}

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "sports.db")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS matches (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			team_a TEXT,
			score_a INTEGER,
			team_b TEXT,
			score_b INTEGER
		);

		CREATE TABLE IF NOT EXISTS teams (
			name TEXT PRIMARY KEY,
			wins INTEGER DEFAULT 0,
			losses INTEGER DEFAULT 0,
			draws INTEGER DEFAULT 0,
			rating REAL DEFAULT 1500
		);
	`)
	if err != nil {
		panic(err)
	}
	return db
}

func GetOrCreateTeam(db *sql.DB, name string) TeamStats {
	var stats TeamStats
	err := db.QueryRow("SELECT wins, losses, draws, rating FROM teams WHERE name = ?", name).Scan(&stats.Wins, &stats.Losses, &stats.Draws, &stats.Rating)
	if err == sql.ErrNoRows {
		_, err = db.Exec("INSERT INTO teams (name) VALUES (?)", name)
		if err != nil {
			panic(err)
		}
		stats.Rating = 1500
		return stats
	} else if err != nil {
		panic(err)
	}
	return stats
}

func UpdateTeamStats(db *sql.DB, name string, win, loss, draw int, newRating float64) {
	_, err := db.Exec(`
		UPDATE teams SET wins = wins + ?, losses = losses + ?, draws = draws + ?, rating = ? WHERE name = ?
	`, win, loss, draw, newRating, name)
	if err != nil {
		panic(err)
	}
}

func RecordMatch(db *sql.DB, m Match) {
	_, err := db.Exec("INSERT INTO matches (team_a, score_a, team_b, score_b) VALUES (?, ?, ?, ?)", m.TeamA, m.ScoreA, m.TeamB, m.ScoreB)
	if err != nil {
		panic(err)
	}

	a := GetOrCreateTeam(db, m.TeamA)
	b := GetOrCreateTeam(db, m.TeamB)

	eA := 1.0 / (1.0 + pow10((b.Rating-a.Rating)/400))
	eB := 1.0 / (1.0 + pow10((a.Rating-b.Rating)/400))
	
	K := 32.0
	var sA, sB float64
	if m.ScoreA > m.ScoreB {
		sA, sB = 1, 0
		UpdateTeamStats(db, m.TeamA, 1, 0, 0, a.Rating+K*(sA-eA))
		UpdateTeamStats(db, m.TeamB, 0, 1, 0, b.Rating+K*(sB-eB))
	} else if m.ScoreA < m.ScoreB {
		sA, sB = 0, 1
		UpdateTeamStats(db, m.TeamA, 0, 1, 0, a.Rating+K*(sA-eA))
		UpdateTeamStats(db, m.TeamB, 1, 0, 0, b.Rating+K*(sB-eB))
	} else {
		sA, sB = 0.5, 0.5
		UpdateTeamStats(db, m.TeamA, 0, 0, 1, a.Rating+K*(sA-eA))
		UpdateTeamStats(db, m.TeamB, 0, 0, 1, b.Rating+K*(sB-eB))
	}
}

func pow10(x float64) float64 {
	return float64(1.0 / (1.0 + pow(10, x)))
}

func pow(base, exp float64) float64 {
	result := 1.0
	for i := 0; i < int(exp); i++ {
		result *= base
	}
	return result
}

func GetTeamStats(db *sql.DB, name string) TeamStats {
	stats := GetOrCreateTeam(db, name)
	return stats
}

func PredictOutcome(db *sql.DB, teamA, teamB string) (float64, float64) {
	a := GetOrCreateTeam(db, teamA)
	b := GetOrCreateTeam(db, teamB)
	eA := 1.0 / (1.0 + pow(10, (b.Rating-a.Rating)/400))
	eB := 1.0 - eA
	return eA, eB
}
