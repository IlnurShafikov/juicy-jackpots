package base

import (
	"github.com/IlnurShafikov/juicy-jackpots/internal/entities/result/events"
	"github.com/IlnurShafikov/juicy-jackpots/pkg/datatransfer/dto/common"
)

// Payload contains the detailed game data for a base spin step.
type Payload struct {
	// Counter indicates the number of free spins or respins triggered (if any)
	Counter *common.Counter `json:"nextCounter,omitempty"`

	// Events contain all game events that occurred (wins, triggers, etc.)
	Events []events.Event `json:"events,omitempty"`

	// Symbols is the board of the spin; coin positions carry the JpToken property
	Symbols []*common.Symbol `json:"symbols"`

	// NextStepType indicates what type of step comes next (base, free, etc.)
	NextStepType int `json:"nextStepType"`
}
