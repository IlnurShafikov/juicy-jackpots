// Package manager provides the main game manager that coordinates game rounds,
// RNG, configuration, and communication between game logic and the API layer.
package manager

import "errors"

// ErrInvalidRTPVersion is returned when a requested RTP version is not configured or supported.
var ErrInvalidRTPVersion = errors.New("invalid rtp version")
