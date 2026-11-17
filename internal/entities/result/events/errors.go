// Package events provides event types for representing game occurrences
// such as wins, triggers, and special features that happen during gameplay.
package events

import "errors"

// ErrPositionNotFound is returned when a symbol position cannot be found in the position map.
// This typically indicates a data inconsistency during event construction.
var ErrPositionNotFound = errors.New("positions not found")
