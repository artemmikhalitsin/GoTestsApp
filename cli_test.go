package poker_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	poker "github.com/artemmikhalitsin/GoTestsApp"
)

var dummyBlindAlerter = &SpyBlindAlerter{}

func TestCLI(t *testing.T) {
	t.Run("Records chris win from cli", func(t *testing.T) {
		in := strings.NewReader("Chris wins")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in, dummyBlindAlerter)
		cli.PlayPoker()

		want := "Chris"

		poker.AssertPlayerWin(t, playerStore, want)
	})

	t.Run("Records cleo win from cli", func(t *testing.T) {
		in := strings.NewReader("Cleo wins")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in, dummyBlindAlerter)
		cli.PlayPoker()

		want := "Cleo"

		poker.AssertPlayerWin(t, playerStore, want)
	})

	t.Run("It schedules printing of blind values", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &poker.StubPlayerStore{}
		blindAlerter := &SpyBlindAlerter{}

		cli := poker.NewCLI(playerStore, in, blindAlerter)
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
