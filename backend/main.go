// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var db = InitDB()

func addMatch(w http.ResponseWriter, r *http.Request) {
	var m Match
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	RecordMatch(db, m)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Match added")
}

func getTeamStatsHandler(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/team/")
	stats := GetTeamStats(db, name)
	json.NewEncoder(w).Encode(stats)
}

func predictHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/predict/"), "/")
	if len(parts) != 2 {
		http.Error(w, "Bad URL format", http.StatusBadRequest)
		return
	}
	teamA, teamB := parts[0], parts[1]
	pA, pB := PredictOutcome(db, teamA, teamB)
	json.NewEncoder(w).Encode(map[string]float64{
		"prob_team_a_win": pA,
		"prob_team_b_win": pB,
	})
}

func main() {
	http.HandleFunc("/match", addMatch)
	http.HandleFunc("/team/", getTeamStatsHandler)
	http.HandleFunc("/predict/", predictHandler)

	fmt.Println("Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
