package poker

import (
	"io"
	"time"
)

// Game is an interface required to run a game
type Game interface {
	Start(numPlayers int, alertsDestination io.Writer)
	Finish(winner string)
}

// TexasHoldem represents a game of texas holdem poker
type TexasHoldem struct {
	store   PlayerStore
	alerter BlindAlerter
}

// NewTexasHoldem creates a new texas holdem game
func NewTexasHoldem(store PlayerStore, alerter BlindAlerter) *TexasHoldem {
	return &TexasHoldem{
		store:   store,
		alerter: alerter,
	}
}

// Start starts the game with a given number of players
func (t *TexasHoldem) Start(numPlayers int, alertsDestination io.Writer) {
	blindIncrement := time.Duration(5+numPlayers) * time.Minute
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		t.alerter.ScheduleAlertAt(blindTime, blind, alertsDestination)
		blindTime += blindIncrement
	}
}

// Finish completes the game with the declared winner
func (t *TexasHoldem) Finish(winner string) {
	t.store.RecordWin(winner)
}
