// Package pairs provides the Pair and Pairs types that represent ordered
// key-value integer mappings, used primarily for jackpot-type multipliers
// and other ID-to-value lookups loaded from CSV configuration files.
package pairs

import "errors"

// ErrInvalidData is returned by [Parse] when a CSV row has fewer than
// the required two columns (key and value).
var ErrInvalidData = errors.New("invalid data")
