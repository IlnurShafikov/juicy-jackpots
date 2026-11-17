// Package datatransfer handles the transformation of internal game results
// into client-facing data structures (DTOs) that are sent via the API.
package datatransfer

import "errors"

// ErrInvalidEventType is returned when an unknown event type is encountered during transformation.
var (
	ErrInvalidEventType = errors.New("invalid event type")
	ErrGameBoardIsNil   = errors.New("game board is nil")
)
