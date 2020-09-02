package poker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"text/template"
)

// PlayerStore stores score information about players
type PlayerStore interface {
	GetPlayerScore(player string) int
	RecordWin(player string)
	GetLeague() League
}

// PlayerServer is a HTTP interface for player information
type PlayerServer struct {
	store PlayerStore
	http.Handler
}

// Player stores info about a player like name and number of wins
type Player struct {
	Name string
	Wins int
}

// NewPlayerServer sets up a new PlayerServer with configured routing
func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := new(PlayerServer)

	p.store = store

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))
	router.Handle("/game", http.HandlerFunc(p.gameHandler))

	p.Handler = router

	return p
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	leagueTable := p.getLeagueTable()

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(leagueTable)
}

func (p *PlayerServer) getLeagueTable() League {
	return p.store.GetLeague()
}

func (p *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")
	switch r.Method {
	case http.MethodGet:
		p.showScore(w, player)
	case http.MethodPost:
		p.processWin(w, player)
	}
}

func (p *PlayerServer) gameHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("game.html")

	if err != nil {
		http.Error(w, fmt.Sprintf("Problem loading template file: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}
