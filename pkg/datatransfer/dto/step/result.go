package step

import "encoding/json"

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

type BetResult struct {
	Steps          []*Step `json:"steps"`
	UserBet        uint64  `json:"userBet"`
	CostMultiplier uint64  `json:"costMultiplier"`
	Seed           int64   `json:"seed"`
}

type LoadGameResult struct {
	GameRules      json.RawMessage `json:"rules"`
	Step           *Step           `json:"step"`
	CostMultiplier uint64          `json:"costMultiplier"`
}
