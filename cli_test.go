package poker_test

import (
	"strings"
	"testing"

	poker "github.com/artemmikhalitsin/GoTestsApp"
)

func TestCLI(t *testing.T) {
	t.Run("Records chris win from cli", func(t *testing.T) {
		in := strings.NewReader("Chris wins")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()

		want := "Chris"

		poker.AssertPlayerWin(t, playerStore, want)
	})

	t.Run("Records cleo win from cli", func(t *testing.T) {
		in := strings.NewReader("Cleo wins")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()

		want := "Cleo"

		poker.AssertPlayerWin(t, playerStore, want)
	})
}
