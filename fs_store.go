package main

import (
	"encoding/json"
	"io"
)

// FileSystemPlayerStore stores leauge results in the file system
type FileSystemPlayerStore struct {
	database io.ReadWriteSeeker
	league   League
}

// NewFileSystemPlayerStore creates a new FileSystemPlayerStore from a given database
func NewFileSystemPlayerStore(database io.ReadWriteSeeker) *FileSystemPlayerStore {
	database.Seek(0, 0)
	league, _ := NewLeague(database)

	return &FileSystemPlayerStore{
		database,
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

	f.database.Seek(0, 0)
	json.NewEncoder(f.database).Encode(f.league)
}
