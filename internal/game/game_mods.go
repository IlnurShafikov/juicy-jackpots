package game

import (
	"encoding/json"
	"fmt"

	"github.com/IlnurShafikov/juicy-jackpots/internal/entities/rules"
)

// RtpGameVersions represents a collection of all game versions with different RTPs.
// Each version contains a set of mods with configurations for the corresponding RTP.
type RtpGameVersions []*Mods

// Mods contains all available game modes for a specific RTP version.
// Includes the configuration for the simple mode, game rules, and RTP version number.
type Mods struct {
	// Simple - configuration for the standard game mode
	Simple *Mod
	// TODO: Add other mods here if needed
	// For example: BuyFeature *Mod for bonus buy mode

	// Rules - JSON with game rules sent to the client
	Rules json.RawMessage

	// RtpVersion - RTP version number (e.g., 96 for 96%)
	RtpVersion int
}

// NewMods creates a new instance of Mods with the specified parameters.
//
// Parameters:
//   - simple: configuration for the simple game mode
//   - rules: JSON representation of game rules
//   - rtpVersion: RTP version number
//
// Returns:
//   - initialized Mods structure
func NewMods(
	simple *Mod, gameRules json.RawMessage, rtpVersion int,
) *Mods {
	return &Mods{
		Simple:     simple,
		Rules:      gameRules,
		RtpVersion: rtpVersion,
	}
}

// GetRTPVersion returns the RTP version for this set of mods.
//
// Returns:
//   - RTP version number (e.g., 96 for 96%)
func (m *Mods) GetRTPVersion() int {
	return m.RtpVersion
}

// GetGameRules returns the game rules in JSON format.
// The rules contain information about payout tables, maximum multipliers, and other parameters.
//
// Returns:
//   - JSON representation of game rules
func (m *Mods) GetGameRules() json.RawMessage {
	return m.Rules
}

// GetGameMod returns a specific game mode by its identifiers.
//
// Parameters:
//   - modId: mode identifier (e.g., "simple" or empty string for base)
//   - featureId: feature identifier (not currently used, leave empty)
//
// Returns:
//   - pointer to the requested mod
//   - ErrFeatureNotImplemented if an unsupported featureId is specified
//   - ErrInvalidGameModId if modId is not found
func (m *Mods) GetGameMod(modId, featureId string) (*Mod, error) {
	// TODO: add feature ID handling if needed
	if featureId != rules.EmptyFeatureID {
		return nil, ErrFeatureNotImplemented
	}

	// TODO: Add other mods here if needed
	switch modId {
	case rules.SimpleModID, rules.EmptyModID:
		return m.Simple, nil
	default:
		return nil, fmt.Errorf("%s: %w", modId, ErrInvalidGameModId)
	}
}
