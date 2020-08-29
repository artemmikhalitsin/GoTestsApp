package main

import (
	"encoding/json"
	"io"
)

// FileSystemPlayerStore stores leauge results in the file system
type FileSystemPlayerStore struct {
	database io.ReadWriteSeeker
}

// GetLeague retrieves the league scores
func (f *FileSystemPlayerStore) GetLeague() League {
	f.database.Seek(0, 0)
	league, _ := NewLeague(f.database)
	return league
}

// GetPlayerScore retrieves a player's score
func (f *FileSystemPlayerStore) GetPlayerScore(name string) (wins int) {
	league := f.GetLeague()
	player := league.Find(name)

	if player != nil {
		return player.Wins
	}

	return 0
}

// RecordWin records a win for a player
func (f *FileSystemPlayerStore) RecordWin(name string) {
	league := f.GetLeague()

	player := league.Find(name)
	player.Wins++

	f.database.Seek(0, 0)
	json.NewEncoder(f.database).Encode(league)
}
