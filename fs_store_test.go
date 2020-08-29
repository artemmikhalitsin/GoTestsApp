package main

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestFSStore(t *testing.T) {
	t.Run("/league from a reader", func(t *testing.T) {

		database, cleanDatabase := createTempFile(t, `[
        {"Name": "Cleo", "Wins": 35},
        {"Name": "Roger", "Wins": 10}
      ]`)
		defer cleanDatabase()
		store := NewFileSystemPlayerStore(database)

		got := store.GetLeague()
		want := League{
			{"Cleo", 35},
			{"Roger", 10},
		}

		assertLeague(t, got, want)

		//read again
		got = store.GetLeague()
		assertLeague(t, got, want)
	})

	t.Run("Get player score", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
        {"Name": "Cleo", "Wins": 35},
        {"Name": "Roger", "Wins": 10}
      ]`)
		defer cleanDatabase()
		store := NewFileSystemPlayerStore(database)

		assertScoreEquals(t, store.GetPlayerScore("Cleo"), 35)
	})

	t.Run("Records a win", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
        {"Name": "Cleo", "Wins": 35},
        {"Name": "Roger", "Wins": 10}
      ]`)
		defer cleanDatabase()
		store := NewFileSystemPlayerStore(database)

		store.RecordWin("Cleo")
		store.RecordWin("Cleo")
		store.RecordWin("Cleo")

		want := 35 + 3
		got := store.GetPlayerScore("Cleo")

		assertScoreEquals(t, got, want)
	})

	t.Run("Record a win for a new player", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
				{"Name": "Cleo", "Wins": 35},
				{"Name": "Roger", "Wins": 10}
			]`)
		defer cleanDatabase()

		newPlayer := "Junior"
		store := NewFileSystemPlayerStore(database)

		store.RecordWin(newPlayer)

		want := 1
		got := store.GetPlayerScore(newPlayer)

		assertScoreEquals(t, got, want)
	})
}

func assertScoreEquals(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("Expected score %d, got score %d", want, got)
	}
}

func createTempFile(t *testing.T, initialData string) (io.ReadWriteSeeker, func()) {
	t.Helper()

	tmpfile, err := ioutil.TempFile("", "db")

	if err != nil {
		t.Errorf("Error creating a temporary file: %v", err)
	}

	tmpfile.Write([]byte(initialData))

	removeFile := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}

	return tmpfile, removeFile
}
