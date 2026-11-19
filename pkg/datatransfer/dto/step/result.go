// Package step contains the DTOs of game steps and bet results sent to the client.
package step

import "encoding/json"

// Step is a single step of a game round as it is sent to the client.
type Step struct {
	Id             []byte `json:"id"`                   // Unique step identifier
	ClientResponse []byte `json:"clientResponse"`       // JSON payload with complete game state
	Cost           uint64 `json:"cost"`                 // How much this spin cost the player
	StepWin        uint64 `json:"stepWin"`              // Win amount for this specific step
	TotalWin       uint64 `json:"totalWin"`             // Cumulative win across all steps
	Type           int    `json:"type"`                 // Step type (TYPE_BASE_SPIN)
	StepsAdded     uint32 `json:"stepsAdded,omitempty"` // Number of bonus steps triggered
	MaxWinReached  bool   `json:"maxWinReached"`        // Whether max win limit was hit
}

// BetResult is the result of a bet: all steps of the round and its metadata.
type BetResult struct {
	// Steps contains all steps of the round in the order they are played.
	Steps []*Step `json:"steps"`
	// UserBet is the bet amount placed by the player.
	UserBet uint64 `json:"userBet"`
	// CostMultiplier is the multiplier applied to the bet to get the round cost.
	CostMultiplier uint64 `json:"costMultiplier"`
	// Seed is the RNG seed of the round (for replay and verification).
	Seed int64 `json:"seed"`
}

// LoadGameResult is the response to the game load request.
type LoadGameResult struct {
	// Step is the initial step shown to the player.
	Step *Step `json:"step"`
	// GameRules contains the game rules as raw JSON.
	GameRules json.RawMessage `json:"rules"`
	// CostMultiplier is the multiplier applied to the bet to get the round cost.
	CostMultiplier uint64 `json:"costMultiplier"`
}
