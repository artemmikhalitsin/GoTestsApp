package poker_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	poker "github.com/artemmikhalitsin/GoTestsApp"
)

var dummyBlindAlerter = &SpyBlindAlerter{}
var dummyPlayerStore = &poker.StubPlayerStore{}
var dummyStdIn = &bytes.Buffer{}
var dummyStdOut = &bytes.Buffer{}

func TestCLI(t *testing.T) {
	t.Run("Start game with 3 players and finish with Mario", func(t *testing.T) {
		numPlayers := 3
		winner := "Mario"
		in := userSends("3", "Mario wins")
		stdout := &bytes.Buffer{}
		game := &SpyGame{}
		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt)
		assertStartedWith(t, game, numPlayers)
		assertFinishedWith(t, game, winner)
	})

	t.Run("It doesn't start the game when given non-numerical input", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		stdin := userSends("Wrong Input")
		game := &SpyGame{}
		cli := poker.NewCLI(stdin, stdout, game)
		cli.PlayPoker()

		assertGameNotStarted(t, game)
		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt, poker.BadInputErrMessage)
	})

	t.Run("It rejects messages that don't contain a win", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		stdin := userSends("3", "Gibberish winner message")
		game := &SpyGame{}
		cli := poker.NewCLI(stdin, stdout, game)
		cli.PlayPoker()

		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt, "Winner message doesn't look right")

		if game.FinishCalled {
			t.Errorf("Game finished with %s, but should not have", game.FinishedWith)
		}
	})
}

type SpyGame struct {
	StartedWith  int
	StartCalled  bool
	FinishedWith string
	FinishCalled bool
}

func (s *SpyGame) Start(numPlayers int, alertsDestination io.Writer) {
	s.StartedWith = numPlayers
	s.StartCalled = true
}

func (s *SpyGame) Finish(winner string) {
	s.FinishedWith = winner
	s.FinishCalled = true
}

func assertScheduledAlert(t *testing.T, got, want scheduledAlert) {
	t.Helper()

	if got.amount != want.amount {
		t.Errorf("got amount %d, expected %d", got.amount, want.amount)
	}

	if got.scheduledAt != want.scheduledAt {
		t.Errorf("got scheule time of %v, expected %v", got.scheduledAt, want.scheduledAt)
	}
}

func assertPrompt(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("Did not get correct prompt. Got: %q, wanted: %q", got, want)
	}
}

func assertMessagesSentToUser(t *testing.T, stdout *bytes.Buffer, messages ...string) {
	t.Helper()
	got := stdout.String()
	want := strings.Join(messages, "")

	if got != want {
		t.Errorf("Got %q sent to stdout, but expected +%v", got, messages)
	}
}

func assertStartedWith(t *testing.T, game *SpyGame, want int) {
	t.Helper()

	if game.StartedWith != want {
		t.Errorf("Started the game with %d players, but got %d", game.StartedWith, want)
	}
}

func assertFinishedWith(t *testing.T, game *SpyGame, want string) {
	t.Helper()

	if game.FinishedWith != want {
		t.Errorf("Got winner %q, but expected %q to win", game.FinishedWith, want)
	}
}

func assertGameNotStarted(t *testing.T, game *SpyGame) {
	t.Helper()

	if game.StartCalled {
		t.Errorf("Game started but it shouldn't have been")
	}
}

func userSends(messages ...string) io.Reader {
	input := strings.Join(messages, "\n")
	return strings.NewReader(input)
}
