package main

import (
	"encoding/json"
	"io"
	"os"
)

// FileSystemPlayerStore stores leauge results in the file system
type FileSystemPlayerStore struct {
	database io.Writer
	league   League
}

// NewFileSystemPlayerStore creates a new FileSystemPlayerStore from a given database
func NewFileSystemPlayerStore(database *os.File) *FileSystemPlayerStore {
	database.Seek(0, 0)
	league, _ := NewLeague(database)

	return &FileSystemPlayerStore{
		&tape{database},
		league,
	}
}

// GetLeague retrieves the league scores
func (f *FileSystemPlayerStore) GetLeague() League {
	return f.league
}

// GetPlayerScore retrieves a player's score
func (f *FileSystemPlayerStore) GetPlayerScore(name string) (wins int) {
	player := f.league.Find(name)

	if player != nil {
		return player.Wins
	}

	return 0
}

// RecordWin records a win for a player
func (f *FileSystemPlayerStore) RecordWin(name string) {
	player := f.league.Find(name)

	if player != nil {
		player.Wins++
	} else {
		f.league = append(f.league, Player{name, 1})
	}

	json.NewEncoder(f.database).Encode(f.league)
}
