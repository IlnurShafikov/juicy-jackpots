// Package rules provides the domain types and constructors for game rules
// that are serialised to JSON and sent to the client on game load.
// Rules include the paytable, volatility level, RTP value, and maximum win multiplier.
package rules

import (
	"encoding/json"
	"fmt"

	jsoniter "github.com/json-iterator/go"

	"github.com/IlnurShafikov/juicy-jackpots/internal/configs"
	"github.com/IlnurShafikov/slot-engine/pkg/lines"
)

// Rules represents the complete game rules sent to the client.
// It includes all information needed for the client to understand the game mechanics and payouts.
type Rules struct {
	// Volatility describes the game's volatility level (risk/variance)
	Volatility *Volatility `json:"volatility"`

	// PayTable contains symbol payouts and match ranges
	PayTable *PayTable `json:"payTable"`

	// JpTypes maps each jackpot type ID (1 MINI .. 4 GRAND) to its prize in units
	JpTypes JpTypes `json:"jpTypes"`

	// Lines contains the paylines: the row of every reel for each line
	Lines lines.Lines `json:"lines"`

	// MaxMultiplier is the maximum win multiplier relative to bet
	MaxMultiplier uint64 `json:"maxMultiplier"`

	// RTP is the expected Return to Player percentage
	RTP float64 `json:"rtp"`
}

// JpTypes maps each jackpot type ID to its prize in units of the cost multiplier.
// It is serialised as a JSON object with the type IDs as string keys.
type JpTypes map[int]uint64

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
	paytable := NewPaytable(general.PayTable)
	jpTypes := makeJackpotMultiplier(general.JpTypes)

	rules := Rules{
		PayTable:      paytable,
		Lines:         general.WinLines,
		JpTypes:       jpTypes,
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

// makeJackpotMultiplier inverts the jackpot types config (prize value → type ID)
// into the JpTypes map (type ID → prize value) expected by the rules JSON payload.
//
// Parameters:
//   - jackpotTypes: map from a jackpot prize value to its jackpot type ID
//
// Returns:
//   - JpTypes map with each jackpot type ID mapped to its prize value
func makeJackpotMultiplier(
	jackpotTypes map[uint64]int,
) JpTypes {
	jpTypes := make(JpTypes, len(jackpotTypes))

	for value, jpType := range jackpotTypes {
		jpTypes[jpType] = value
	}

	return jpTypes
}
