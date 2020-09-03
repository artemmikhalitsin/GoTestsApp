package poker_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"testing"
	"time"

	poker "github.com/artemmikhalitsin/GoTestsApp"
)

func TestGame_Start(t *testing.T) {
	t.Run("It schedules blind alerts for 5 players", func(t *testing.T) {
		playerStore := &poker.StubPlayerStore{}
		blindAlerter := &SpyBlindAlerter{}
		game := poker.NewTexasHoldem(playerStore, blindAlerter)

		game.Start(5, ioutil.Discard)

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

		assertAlertsScheduled(t, blindAlerter.alerts, cases)
	})

	t.Run("It schedules blind alerts for 7 players", func(t *testing.T) {
		numPlayers := 7
		blindAlerter := &SpyBlindAlerter{}
		game := poker.NewTexasHoldem(dummyPlayerStore, blindAlerter)

		game.Start(numPlayers, ioutil.Discard)

		cases := []scheduledAlert{
			{0 * time.Second, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}

		assertAlertsScheduled(t, blindAlerter.alerts, cases)
	})
}

func TestGame_Finish(t *testing.T) {
	t.Run("Records cleo win from cli", func(t *testing.T) {
		winner := "Cleo"
		playerStore := &poker.StubPlayerStore{}
		game := poker.NewTexasHoldem(playerStore, dummyBlindAlerter)

		game.Finish(winner)

		poker.AssertPlayerWin(t, playerStore, winner)
	})
}

type SpyBlindAlerter struct {
	alerts []scheduledAlert
}

type scheduledAlert struct {
	scheduledAt time.Duration
	amount      int
}

func (s *SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int, to io.Writer) {
	s.alerts = append(s.alerts, scheduledAlert{duration, amount})
}

func assertAlertsScheduled(t *testing.T, got []scheduledAlert, want []scheduledAlert) {
	for i, alert := range want {
		t.Run(fmt.Sprint(alert), func(t *testing.T) {
			if len(got) <= i {
				t.Fatalf("alert %d was not scheduled. %v", i, got)
			}

			scheduled := got[i]

			assertScheduledAlert(t, scheduled, alert)
		})
	}
}
