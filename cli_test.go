package poker_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	poker "github.com/artemmikhalitsin/GoTestsApp"
)

var dummyBlindAlerter = &SpyBlindAlerter{}
var dummyPlayerStore = &poker.StubPlayerStore{}
var dummyStdIn = &bytes.Buffer{}
var dummyStdOut = &bytes.Buffer{}

func TestCLI(t *testing.T) {
	t.Run("Records chris win from cli", func(t *testing.T) {
		in := strings.NewReader("7\nChris wins")
		playerStore := &poker.StubPlayerStore{}
		game := poker.NewGame(playerStore, dummyBlindAlerter)

		cli := poker.NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()

		want := "Chris"

		poker.AssertPlayerWin(t, playerStore, want)
	})

	t.Run("Records cleo win from cli", func(t *testing.T) {
		in := strings.NewReader("7\nCleo wins")
		playerStore := &poker.StubPlayerStore{}
		game := poker.NewGame(playerStore, dummyBlindAlerter)

		cli := poker.NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()

		want := "Cleo"

		poker.AssertPlayerWin(t, playerStore, want)
	})

	t.Run("It schedules printing of blind values", func(t *testing.T) {
		in := strings.NewReader("5\nChris wins\n")
		playerStore := &poker.StubPlayerStore{}
		blindAlerter := &SpyBlindAlerter{}
		game := poker.NewGame(playerStore, blindAlerter)

		cli := poker.NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()

		cases := []scheduledAlert{
			{0 * time.Second, 100},
			{10 * time.Minute, 200},
			{20 * time.Minute, 300},
			{30 * time.Minute, 400},
			{40 * time.Minute, 500},
			{50 * time.Minute, 600},
			{60 * time.Minute, 800},
			{70 * time.Minute, 1000},
			{80 * time.Minute, 2000},
			{90 * time.Minute, 4000},
			{100 * time.Minute, 8000},
		}

		for i, c := range cases {
			t.Run(fmt.Sprintf("%d scheduled for %v", c.amount, c.scheduledAt), func(t *testing.T) {
				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}
				alert := blindAlerter.alerts[i]
				assertScheduledAlert(t, alert, c)
			})
		}
	})

	t.Run("It prompts for the number of players", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		game := poker.NewGame(dummyPlayerStore, dummyBlindAlerter)
		cli := poker.NewCLI(dummyStdIn, stdout, game)
		cli.PlayPoker()

		assertPrompt(t, stdout.String(), poker.PlayerPrompt)
	})

	t.Run("It provides correct blinds given a number of players", func(t *testing.T) {
		in := strings.NewReader("7\n")
		blindAlerter := &SpyBlindAlerter{}
		game := poker.NewGame(dummyPlayerStore, blindAlerter)
		cli := poker.NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()

		cases := []scheduledAlert{
			{0 * time.Second, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}

		for i, want := range cases {
			t.Run(fmt.Sprint(want), func(t *testing.T) {
				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled. %v", i, blindAlerter.alerts)
				}

				got := blindAlerter.alerts[i]

				assertScheduledAlert(t, got, want)
			})
		}
	})
}

type SpyBlindAlerter struct {
	alerts []scheduledAlert
}

type scheduledAlert struct {
	scheduledAt time.Duration
	amount      int
}

func (s *SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int) {
	s.alerts = append(s.alerts, scheduledAlert{duration, amount})
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
	if got != want {
		t.Errorf("Did not get correct prompt. Got: %q, wanted: %q", got, poker.PlayerPrompt)
	}
}
