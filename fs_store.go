package main

import (
	"io"
)

// FileSystemPlayerStore stores leauge results in the file system
type FileSystemPlayerStore struct {
	database io.ReadSeeker
}

// GetLeague retrieves the league scores
func (f *FileSystemPlayerStore) GetLeague() []Player {
	f.database.Seek(0, 0)
	league, _ := NewLeague(f.database)
	return league
}

// GetPlayerScore retrieves a player's score
func (f *FileSystemPlayerStore) GetPlayerScore(name string) (wins int) {
	for _, player := range f.GetLeague() {
		if player.Name == name {
			wins = player.Wins
			return
		}
	}
	return
}
