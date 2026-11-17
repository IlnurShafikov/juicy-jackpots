// Package collectors provides components for collecting and managing game events
// such as wins, multipliers, and other game mechanics during gameplay.
package collectors

import "errors"

// ErrInvalidMaxWinLimit is returned when the maximum win limit is invalid (must be greater than 0).
var ErrInvalidMaxWinLimit = errors.New("maxWin must be greater than 0")
