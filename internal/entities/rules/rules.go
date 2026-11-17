// Package rules provides the domain types and constructors for game rules
// that are serialised to JSON and sent to the client on game load.
// Rules include the paytable, volatility level, RTP value, and maximum win multiplier.
package rules

import (
	"encoding/json"
	"fmt"

	"github.com/IlnurShafikov/juicy-jackpots/internal/configs"
	"github.com/IlnurShafikov/slot-engine/pkg/paytable"
	jsoniter "github.com/json-iterator/go"
)

// Rules represents the complete game rules sent to the client.
// It includes all information needed for the client to understand the game mechanics and payouts.
// TODO: implement the required fields for your specific game
type Rules struct {
	// Volatility describes the game's volatility level (risk/variance)
	Volatility *Volatility `json:"volatility"`

	// PayTable contains symbol payouts and match ranges
	PayTable *paytable.PayTable `json:"payTable"`

	// MaxMultiplier is the maximum win multiplier relative to bet
	MaxMultiplier uint64 `json:"maxMultiplier"`

	// RTP is the expected Return to Player percentage
	RTP float64 `json:"rtp"`
}

// MakeGameRules constructs the game rules JSON from configuration data.
// This JSON is sent to the client to inform them of the game's mechanics and payouts.
//
// Parameters:
//   - general: general configuration containing paytable and multipliers
//   - expRTP: expected RTP percentage for this game version
//
// Returns:
//   - json.RawMessage containing the serialized game rules
//   - error if marshaling fails
func MakeGameRules(
	general *configs.GeneralConfig,
	expRTP float64,
) (json.RawMessage, error) {
	rules := Rules{
		PayTable:      general.PayTable,
		MaxMultiplier: general.Multipliers.MaxWinMultiplier,
		RTP:           expRTP,
		Volatility:    &Volatility{Level: general.Volatility},
	}

	rulesJSON, err := jsoniter.Marshal(rules)
	if err != nil {
		return nil, fmt.Errorf("marshal rules: %w", err)
	}

	return rulesJSON, nil
}

// Volatility represents the game's volatility (variance/risk) level.
// Higher volatility means larger but less frequent wins,
// lower volatility means smaller but more frequent wins.
type Volatility struct {
	// Level is the volatility classification (e.g., "low", "medium", "high")
	Level string `json:"level"`
}
