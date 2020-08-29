package main

import (
	"strings"
	"testing"
)

func TestFSStore(t *testing.T) {
	t.Run("/league from a reader", func(t *testing.T) {
		database := strings.NewReader(`[
        {"Name": "Cleo", "Wins": 35},
        {"Name": "Roger", "Wins": 10}
      ]`)

		store := FileSystemPlayerStore{database}

		got := store.GetLeague()
		want := []Player{
			{"Cleo", 35},
			{"Roger", 10},
		}

		assertLeague(t, got, want)

		//read again
		got = store.GetLeague()
		assertLeague(t, got, want)
	})

	t.Run("Get player score", func(t *testing.T) {
		database := strings.NewReader(`[
        {"Name": "Cleo", "Wins": 35},
        {"Name": "Roger", "Wins": 10}
      ]`)
		player := "Cleo"

		store := FileSystemPlayerStore{database}

		got := store.GetPlayerScore(player)
		want := 35

		if got != want {
			t.Errorf("Expected score %d, got score %d", want, got)
		}
	})
}
