package main

import "sync"

// InMemoryPlayerStore collects data about players and stores in memory
type InMemoryPlayerStore struct {
	store map[string]int
	mu    sync.RWMutex
}

// NewInMemoryPlayerStore initialzes a new player store
func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{
		map[string]int{},
		sync.RWMutex{},
	}
}

// GetPlayerScore retrieves a player's score from the store
func (i *InMemoryPlayerStore) GetPlayerScore(player string) int {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.store[player]
}

// RecordWin records a win for a player
func (i *InMemoryPlayerStore) RecordWin(player string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.store[player]++
}

// GetLeague retrieves the players currently in the league
func (i *InMemoryPlayerStore) GetLeague() []Player {
	var league []Player

	for name, wins := range i.store {
		league = append(league, Player{name, wins})
	}

	return league
}
