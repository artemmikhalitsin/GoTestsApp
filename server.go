package poker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

const htmlTemplatePath = "game.html"

// PlayerStore stores score information about players
type PlayerStore interface {
	GetPlayerScore(player string) int
	RecordWin(player string)
	GetLeague() League
}

// PlayerServer is a HTTP interface for player information
type PlayerServer struct {
	store PlayerStore
	game  Game
	http.Handler
	template *template.Template
}

// Player stores info about a player like name and number of wins
type Player struct {
	Name string
	Wins int
}

// NewPlayerServer sets up a new PlayerServer with configured routing
func NewPlayerServer(store PlayerStore, game Game) (*PlayerServer, error) {
	p := new(PlayerServer)

	tmpl, err := template.ParseFiles(htmlTemplatePath)

	if err != nil {
		return nil, fmt.Errorf("Problem loading template file %s: %v", htmlTemplatePath, err)
	}

	p.template = tmpl
	p.store = store
	p.game = game

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))
	router.Handle("/game", http.HandlerFunc(p.gameHandler))
	router.Handle("/ws", http.HandlerFunc(p.websocketHandler))

	p.Handler = router

	return p, nil
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
	p.template.Execute(w, nil)
}

func (p *PlayerServer) websocketHandler(w http.ResponseWriter, r *http.Request) {
	ws := newPlayerServerWS(w, r)
	defer ws.Conn.Close()

	playersMsg := ws.WaitForMsg()
	numPlayers, _ := strconv.Atoi(string(playersMsg))
	p.game.Start(numPlayers, ws)

	winner := ws.WaitForMsg()
	p.game.Finish(string(winner))
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
