// Package common contains common data transfer objects (DTOs) used across different game phases.
// These structures are serialized and sent to the client.
package common

// Counter tracks the state of a feature that has multiple steps (free spins, respins).
// It provides information about how many steps were added,
// how many have been played, and how many remain.
type Counter struct {
	// Added is the number of steps added in the current step
	Added uint32 `json:"added"`

	// Left is the number of steps remaining to be played
	Left uint32 `json:"left"`

	// Played is the number of steps that have been completed
	Played uint32 `json:"played"`
}
