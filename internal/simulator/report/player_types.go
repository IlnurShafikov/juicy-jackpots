//nolint:all // it's ok for simulator
package report

import (
	"sort"
)

// Player represents a single player with balance and spin history
type Player struct {
	InitialBalance uint64
	CurrentBalance uint64
	SpinsCount     int
}

// NewPlayer creates a new player with initial balance
func NewPlayer(balance uint64) *Player {
	return &Player{
		InitialBalance: balance,
		CurrentBalance: balance,
		SpinsCount:     0,
	}
}

// CanSpin checks if player has enough balance to make a spin
func (p *Player) CanSpin(betCost uint64) bool {
	return p.CurrentBalance >= betCost
}

// MakeBet deducts bet cost from player's balance
func (p *Player) MakeBet(betCost uint64) {
	p.CurrentBalance -= betCost
}

// AddWin adds win amount to player's balance
func (p *Player) AddWin(winAmount uint64) {
	p.CurrentBalance += winAmount
}

// IncrementSpins increments the number of spins made
func (p *Player) IncrementSpins() {
	p.SpinsCount++
}

// PlayerGroup represents a group of players with same initial balance
type PlayerGroup struct {
	InitialBalance uint64
	Players        []*Player
	PlayersCount   int
}

// NewPlayerGroup creates a new group of players
func NewPlayerGroup(balance uint64, playersCount int) *PlayerGroup {
	players := make([]*Player, playersCount)

	for i := range playersCount {
		players[i] = NewPlayer(balance)
	}

	return &PlayerGroup{
		InitialBalance: balance,
		Players:        players,
		PlayersCount:   playersCount,
	}
}

// PlayerStats contains statistical results for a player group
type PlayerStats struct {
	InitialBalance uint64
	PlayersCount   int
	MedianSpins    float64
	AllSpinCounts  []int
}

// NewPlayerStats creates stats from a player group
func NewPlayerStats(group *PlayerGroup) *PlayerStats {
	spinCounts := make([]int, len(group.Players))
	for i, player := range group.Players {
		spinCounts[i] = player.SpinsCount
	}

	return &PlayerStats{
		InitialBalance: group.InitialBalance,
		PlayersCount:   group.PlayersCount,
		MedianSpins:    calculateMedian(spinCounts),
		AllSpinCounts:  spinCounts,
	}
}

// calculateMedian calculates median value from a slice of integers
func calculateMedian(values []int) float64 {
	if len(values) == 0 {
		return 0
	}

	resSorted := make([]int, len(values))
	copy(resSorted, values)
	sort.Ints(resSorted)

	n := len(resSorted)
	if n%2 == 0 {
		return float64(resSorted[n/2-1]+resSorted[n/2]) / 2.0
	}

	return float64(resSorted[n/2])
}
