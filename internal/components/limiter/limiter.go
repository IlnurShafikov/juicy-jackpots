// Package limiter provides win limiting functionality to cap maximum payouts.
package limiter

import (
	"fmt"

	"github.com/IlnurShafikov/slot-engine/pkg/tools/maxwin"
)

// WinLimiter defines the interface for limiting and tracking win amounts.
// Implementations ensure that total wins do not exceed configured limits.
type WinLimiter interface {
	// Add attempts to add a payout to the total win, applying limits if necessary
	Add(payout uint64) (limitedPayout uint64, maxWinReached bool)

	// MaxWinReached returns true if the maximum win limit has been reached
	MaxWinReached() bool

	// TotalWin returns the cumulative total win amount
	TotalWin() uint64
}

var _ WinLimiter = (*maxwin.WinLimiter)(nil)

// New creates a new win limiter with the specified maximum win limit.
//
// Parameters:
//   - maxWinLimit: the maximum allowed win amount
//
// Returns:
//   - pointer to a new maxwin.Collector
//   - error if the limiter creation fails
func New(maxWinLimit uint64) (*maxwin.WinLimiter, error) {
	resp, err := maxwin.NewWinLimiter(maxWinLimit)
	if err != nil {
		return nil, fmt.Errorf("creating win limiter maxWinLimit=[%d]: %w",
			maxWinLimit, err)
	}

	return resp, nil
}
