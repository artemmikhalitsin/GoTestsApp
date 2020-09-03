package poker

import (
	"io"
)

// Game is an interface required to run a game
type Game interface {
	Start(numPlayers int, alertsDestination io.Writer)
	Finish(winner string)
}
