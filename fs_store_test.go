package poker_test

import (
	"io/ioutil"
	"os"
	"testing"

	poker "github.com/artemmikhalitsin/GoTestsApp"
)

func TestFSStore(t *testing.T) {
	t.Run("/league from a reader is in order from highest score to lowest", func(t *testing.T) {

		database, cleanDatabase := createTempFile(t, `[
        {"Name": "Cleo", "Wins": 10},
        {"Name": "Roger", "Wins": 35},
				{"Name": "Cedar", "Wins": 22}
      ]`)
		defer cleanDatabase()
		store, err := poker.NewFileSystemPlayerStore(database)

		assertNoError(t, err)

		got := store.GetLeague()
		// ordered from highest score to lowest
		want := poker.League{
			{"Roger", 35},
			{"Cedar", 22},
			{"Cleo", 10},
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
		store, err := poker.NewFileSystemPlayerStore(database)

		assertNoError(t, err)
		assertScoreEquals(t, store.GetPlayerScore("Cleo"), 35)
	})

	t.Run("Records a win", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
        {"Name": "Cleo", "Wins": 35},
        {"Name": "Roger", "Wins": 10}
      ]`)
		defer cleanDatabase()
		store, err := poker.NewFileSystemPlayerStore(database)

		assertNoError(t, err)

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
		store, err := poker.NewFileSystemPlayerStore(database)

		assertNoError(t, err)

		store.RecordWin(newPlayer)

		want := 1
		got := store.GetPlayerScore(newPlayer)

		assertScoreEquals(t, got, want)
	})

	t.Run("Works with an empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()

		_, err := poker.NewFileSystemPlayerStore(database)

		assertNoError(t, err)
	})
}

func assertScoreEquals(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("Expected score %d, got score %d", want, got)
	}
}

func createTempFile(t *testing.T, initialData string) (*os.File, func()) {
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
