// Package coin defines the Result type for jackpot hunt coin placement.
package coin

import "github.com/IlnurShafikov/slot-engine/pkg/model"

// Result holds the outcome of the jackpot hunt coin placement for a single spin.
type Result struct {
	// Positions lists the board cells where coins were placed this spin.
	Positions model.Positions `json:"positions"`
	// Value is the jackpot payout value associated with the coin count drawn this spin.
	Value uint64 `json:"value"`
	// JpType is the jackpot type identifier that corresponds to Value.
	JpType int `json:"jackpotType"`
}

// NewResult creates a new coin placement result with zero-initialized fields.
//
// Returns:
//   - pointer to a new Result instance ready to be populated
func NewResult() *Result {
	return &Result{
		Positions: nil,
		Value:     0,
		JpType:    0,
	}
}

// Copy returns a shallow copy of the Result.
// It duplicates all fields but does not deep-copy the Positions slice.
//
// Returns:
//   - pointer to a new Result with the same field values
func (r Result) Copy() *Result {
	return &Result{
		Positions: r.Positions.Copy(),
		Value:     r.Value,
		JpType:    r.JpType,
	}
}
